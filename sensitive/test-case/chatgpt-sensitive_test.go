package testcase

import (
	"chatgpt-sensitive/pkg/config"
	"chatgpt-sensitive/pkg/filter"
	"chatgpt-sensitive/proto"
	"chatgpt-sensitive/server"
	"net"

	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"os"
	"testing"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

var configPath = flag.String("config", "../dev.config.ymal", "单元测试配置文件")
var dictPath = flag.String("dict", "./dict.txt", "敏感词库地址")

// go test ./test-case/. --config=../dev.config.yaml --dict=./dict.txt
func TestMain(m *testing.M) {
	flag.Parse()
	config.InitConf(*configPath)
	filter.InitFilter(*dictPath)
	m.Run()
}

func TestValidate(t *testing.T) {
	dataList := []struct {
		req *proto.ValidateReq
		res *proto.ValidateRes
	}{
		{
			req: &proto.ValidateReq{
				Text: "PKfadf TEL niuiug",
			},
			res: &proto.ValidateRes{
				Ok:   false,
				Word: "TEL",
			},
		},
		{
			req: &proto.ValidateReq{
				Text: "请帮我写一篇1000字的计算机论文好嘛",
			},
			res: &proto.ValidateRes{
				Ok:   false,
				Word: "论文",
			},
		},
		{
			req: &proto.ValidateReq{
				Text: "请帮我写一hello你好文好嘛",
			},
			res: &proto.ValidateRes{
				Ok:   false,
				Word: "hello你好",
			},
		},
		{
			req: &proto.ValidateReq{
				Text: "今天天气怎么样",
			},
			res: &proto.ValidateRes{
				Ok:   true,
				Word: "",
			},
		},
	}

	// server
	lis, err := net.Listen("tcp", "localhost:50053")
	if err != nil {
		t.Error(err)
		return
	}
	// 添加服务器启动选项
	s := grpc.NewServer(getServerOptions()...)
	defer s.Stop()
	sensitiveWordServer := server.NewSensitiveWordServer(filter.GetFilter())
	proto.RegisterChatGPTSensitiveServer(s, sensitiveWordServer)
	go func() {
		err = s.Serve(lis)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	// client
	conn, err := grpc.Dial("localhost:50053", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	cli := proto.NewChatGPTSensitiveClient(conn)
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	for _, item := range dataList {
		res, err := cli.Validate(ctx, item.req)
		if err != nil {
			t.Error(err)
		}
		if res.Ok != item.res.Ok || res.Word != item.res.Word {
			t.Error("敏感词过滤结果与预期不一致")
		}
	}
}

func getClientOptions() []grpc.DialOption {
	clientOption := make([]grpc.DialOption, 0)
	clientOption = append(clientOption, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return clientOption
}

func getTlsOpt(cert, serviceName string) grpc.DialOption {
	creds, err := credentials.NewClientTLSFromFile(cert, serviceName)
	if err != nil {
		panic(err)
	}
	return grpc.WithTransportCredentials(creds)
}

func getMTLSOpt(caCert, certFile, keyFile string) grpc.DialOption {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	ca := x509.NewCertPool()
	bytes, err := os.ReadFile(caCert)
	if err != nil {
		panic(err)
	}
	ok := ca.AppendCertsFromPEM(bytes)
	if !ok {
		panic("append cert failed")
	}
	tlsConfig := &tls.Config{
		ServerName:   "chatgpt-data.grpc.lin.com",
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}
	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
}

func getAuth() grpc.DialOption {
	token := config.GetConf().Server.AccessToken
	perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
	return grpc.WithPerRPCCredentials(perRPC)
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
