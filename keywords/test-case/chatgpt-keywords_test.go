package testcase

import (
	"chatgpt-keywords/pkg/config"
	"chatgpt-keywords/pkg/filter"
	"chatgpt-keywords/proto"
	"chatgpt-keywords/server"
	"fmt"
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
var dictPath = flag.String("dict", "../cainiao-coding.txt", "关键词库地址")

// go test ./test-case/. --config=../dev.config.yaml --dict=../cainiao-coding.txt
func TestMain(m *testing.M) {
	flag.Parse()
	config.InitConf(*configPath)
	filter.InitFilter(*dictPath)
	m.Run()
}

func TestFindAll(t *testing.T) {
	dataList := []struct {
		req *proto.FindAllReq
		res *proto.FindAllRes
	}{
		{
			req: &proto.FindAllReq{
				Text: "aninpog access activate parameter fmoienqgor",
			},
			res: &proto.FindAllRes{
				Words: []string{"access", "activate", "parameter"},
			},
		},
		{
			req: &proto.FindAllReq{
				Text: "access activate parameter",
			},
			res: &proto.FindAllRes{
				Words: []string{"access", "activate", "parameter"},
			},
		},
		{
			req: &proto.FindAllReq{
				Text: "access指针协议prototype",
			},
			res: &proto.FindAllRes{
				Words: []string{"指针", "协议"},
			},
		},
		{
			req: &proto.FindAllReq{
				Text: "你可能需要Common Lisp 对象系统是吧",
			},
			res: &proto.FindAllRes{
				Words: []string{"Common Lisp 对象系统", "对象"},
			},
		},
	}

	// server
	lis, err := net.Listen("tcp", "localhost:50054")
	if err != nil {
		t.Error(err)
		return
	}
	// 添加服务器启动选项
	s := grpc.NewServer(getServerOptions()...)
	defer s.Stop()
	keyWordServer := server.NewKeyWordsServer(filter.GetFilter())
	proto.RegisterChatGPTKeywordsServer(s, keyWordServer)
	go func() {
		err = s.Serve(lis)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	// client
	conn, err := grpc.Dial("localhost:50054", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	cli := proto.NewChatGPTKeywordsClient(conn)
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	for _, item := range dataList {
		res, err := cli.FindAll(ctx, item.req)
		if err != nil {
			t.Error(err)
		}
		if len(res.Words) != len(item.res.Words) {
			fmt.Println(res.Words)
			fmt.Println(item.res.Words)
			t.Error("关键词过滤结果与预期不一致")
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
