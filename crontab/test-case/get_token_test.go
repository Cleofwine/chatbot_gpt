package casetest

import (
	"chatgpt-crontab/pkg/config"
	"chatgpt-crontab/pkg/db/redis"
	"chatgpt-crontab/pkg/log"
	"chatgpt-crontab/proto"
	"chatgpt-crontab/token-server/server"
	"context"
	"flag"
	"fmt"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var configPath = flag.String("config", "../dev.config.ymal", "单元测试配置文件")

// go test ./test-case/. --config=../dev.config.yaml -v
func TestMain(m *testing.M) {
	flag.Parse()
	config.InitConf(*configPath)
	redis.InitRedisPool()
	go startGRPCServer()
	m.Run()
}

func startGRPCServer() {
	cnf := config.GetConf()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 50056))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(getServerOptions()...)
	logger := log.NewLogger()
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)
	tokenServer := server.NewTokenServer(cnf, logger)
	proto.RegisterTokenServer(s, tokenServer)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func getClientOptions() []grpc.DialOption {
	clientOption := make([]grpc.DialOption, 0)
	clientOption = append(clientOption, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return clientOption
}

func appendBearerTokenToContext(ctx context.Context) context.Context {
	token := config.GetConf().Server.AccessToken
	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(ctx, md)
}

func getServerOptions() []grpc.ServerOption {
	var opts = make([]grpc.ServerOption, 0)
	opts = append(opts, server.GetKeepaliveOpt()...)
	opts = append(opts, grpc.StreamInterceptor(server.StreamInterceptor))
	opts = append(opts, grpc.UnaryInterceptor(server.UnaryInterceptor))
	return opts
}

func TestGetWxOfficialToken(t *testing.T) {
	conn, err := grpc.Dial("localhost:50056", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	client := proto.NewTokenClient(conn)
	in := &proto.TokenRequest{
		Id:  config.GetConf().WxOfficials[0].AppId,
		App: "",
		Typ: proto.TokenType_WECHATOFFICIAL,
	}
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	res, err := client.GetToken(ctx, in)
	if err != nil {
		t.Error(err)
		return
	}
	if res.AccessToken == "" {
		t.Error("access_token获取失败")
		return
	}
	t.Log(res.AccessToken)
}

func TestGetWeComToken(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:50056", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	client := proto.NewTokenClient(conn)
	in := &proto.TokenRequest{
		Id:  config.GetConf().WeComs[0].CorpId,
		App: config.GetConf().WeComs[0].App,
		Typ: proto.TokenType_WECOM,
	}
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	res, err := client.GetToken(ctx, in)
	if err != nil {
		t.Error(err)
		return
	}
	if res.AccessToken == "" {
		t.Error("access_token获取失败")
		return
	}
	t.Log(res.AccessToken)
}
