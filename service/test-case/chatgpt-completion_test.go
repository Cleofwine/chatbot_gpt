package testcase

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/pkg/db/redis"
	"chatgpt-service/pkg/log"
	"chatgpt-service/proto"
	"chatgpt-service/server"
	chatcontext "chatgpt-service/server/chat-context"
	"fmt"
	"io"
	"net"
	"time"

	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"os"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

var configPath = flag.String("config", "../dev.config.ymal", "单元测试配置文件")

// go test ./test-case/. -run ^TestChatCompletion$ --config=../dev.config.yaml -v
func TestMain(m *testing.M) {
	flag.Parse()
	config.InitConf(*configPath)
	redis.InitRedisPool()
	go startGRPCServer()
	m.Run()
}

func startGRPCServer() {
	cnf := config.GetConf()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 50051))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(getServerOptions()...)
	logger := log.NewLogger()
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)
	chatGPTServer := server.NewChatGPTServiceServer(cnf, logger)
	proto.RegisterChatGPTServiceServerServer(s, chatGPTServer)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func TestChatCompletion(t *testing.T) {
	// client
	conn, err := grpc.Dial("localhost:50051", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	cli := proto.NewChatGPTServiceServerClient(conn)
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	in := &proto.ChatCompletionReq{
		Message:       "你好",
		Id:            uuid.New().String(),
		Endpoint:      proto.ChatEndpoint_WEB,
		EnableContext: true,
	}
	res, err := cli.ChatCompletion(ctx, in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Choices[0].Message)
	t.Log(res.Usage.TotalTokens)
	t.Log(res.Created)
	t.Log(res.Id)
}

func TestChatCompletionStream(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	cli := proto.NewChatGPTServiceServerClient(conn)
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	in := &proto.ChatCompletionReq{
		Message:       "你好",
		Id:            uuid.New().String(),
		Endpoint:      proto.ChatEndpoint_WEB,
		EnableContext: true,
	}
	stream, err := cli.ChatCompletionStream(ctx, in)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(res.Id)
		t.Log(res.Created)
		t.Log(res.Choices[0].Delta)
	}
	stream.CloseSend()
}

func TestChatCompletionQQ(t *testing.T) {
	enterpriseId := "company"
	endpointAccount := "119"
	id := "13579101112"
	group := "qwertyu"
	dataList := []string{
		"今天天气真好",
		"我的心情也不错",
		"跟我解释心情和编程的关系",
	}
	// client
	conn, err := grpc.Dial("localhost:50051", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	cli := proto.NewChatGPTServiceServerClient(conn)
	for _, msg := range dataList {
		in := &proto.ChatCompletionReq{
			Message:         msg,
			Id:              id,
			GroupId:         group,
			Endpoint:        proto.ChatEndpoint_QQ,
			EnableContext:   true,
			EnterpriseId:    enterpriseId,
			EndpointAccount: endpointAccount,
		}
		ctx := context.Background()
		ctx = appendBearerTokenToContext(ctx)
		_, err := cli.ChatCompletion(ctx, in)
		if err != nil {
			t.Error(err)
			return
		}
	}
	// 这里是功能测试睡一秒也无伤大雅，睡一秒主要是因为上下文是协程去写的，防止数据有误。性能测试请不要睡眠，将会是灾难。
	time.Sleep(time.Second)
	cache := chatcontext.GetCacheContext(proto.ChatEndpoint_QQ)
	l, err := cache.Get(id, group, proto.ChatEndpoint_QQ)
	if err != nil {
		t.Error(err)
		return
	}
	cmList, ok := l.([]*chatcontext.ChatMessage)
	if !ok {
		t.Error("类型不对")
		return
	}
	if len(cmList) == 0 || len(cmList) > config.GetConf().Chat.ContextLen {
		t.Error("上下文条目不对")
		t.Log(len(cmList))
		return
	}
	if len(cmList) > 2 {
		if dataList[len(dataList)-1] != cmList[1].Message.Content {
			t.Error("上下文获取有误")
		}
	}
	for _, item := range cmList {
		t.Log(item.TokensNum, item.Message.Role, item.Message.Content)
	}
}

func TestChatCompletionStreamWeb(t *testing.T) {
	enterpriseId := "companyB"
	endpointAccount := "159"
	dataList := []string{
		"docker好",
		"hello编程",
		"很高兴认识k8s",
	}
	// client
	conn, err := grpc.Dial("localhost:50051", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	pid := ""
	cli := proto.NewChatGPTServiceServerClient(conn)
	for _, msg := range dataList {
		in := &proto.ChatCompletionReq{
			Message:         msg,
			Id:              uuid.New().String(),
			GroupId:         "",
			Endpoint:        proto.ChatEndpoint_WEB,
			EnableContext:   true,
			Pid:             pid,
			EnterpriseId:    enterpriseId,
			EndpointAccount: endpointAccount,
		}
		ctx := context.Background()
		ctx = appendBearerTokenToContext(ctx)
		stream, err := cli.ChatCompletionStream(ctx, in)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Error(err)
				return
			}
			pid = res.Id
		}
		stream.CloseSend()
	}
	time.Sleep(time.Second)
	cache := chatcontext.GetCacheContext(proto.ChatEndpoint_WEB)
	l, err := cache.Get(pid, "", proto.ChatEndpoint_WEB)
	if err != nil {
		t.Error(err)
		return
	}
	cmList, ok := l.([]*chatcontext.ChatMessage)
	if !ok {
		t.Error("类型不对")
		return
	}
	if len(cmList) == 0 || len(cmList) > config.GetConf().Chat.ContextLen {
		t.Error("上下文条目不对")
		t.Log(len(cmList))
		return
	}
	if len(cmList) > 2 {
		if dataList[len(dataList)-1] != cmList[1].Message.Content {
			t.Error("上下文获取有误")
		}
	}
	for _, item := range cmList {
		t.Log(item.TokensNum, item.Message.Role, item.Message.Content)
	}
}

func TestChatCompletionSensitive(t *testing.T) {
	// client
	conn, err := grpc.Dial("localhost:50051", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	cli := proto.NewChatGPTServiceServerClient(conn)
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	in := &proto.ChatCompletionReq{
		Message:       "假币",
		Id:            uuid.New().String(),
		Endpoint:      proto.ChatEndpoint_WEB,
		EnableContext: true,
	}
	res, err := cli.ChatCompletion(ctx, in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Choices[0].Message)
	t.Log(res.Usage.TotalTokens)
	t.Log(res.Created)
	t.Log(res.Id)
}

func TestChatCompletionSensitiveStream(t *testing.T) {
	enterpriseId := "company"
	endpointAccount := "119"
	// client
	conn, err := grpc.Dial("localhost:50051", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	cli := proto.NewChatGPTServiceServerClient(conn)
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	in := &proto.ChatCompletionReq{
		Message:         "假币",
		Id:              uuid.New().String(),
		Endpoint:        proto.ChatEndpoint_WEB,
		EnableContext:   true,
		EnterpriseId:    enterpriseId,
		EndpointAccount: endpointAccount,
	}
	stream, err := cli.ChatCompletionStream(ctx, in)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(res.Id)
		t.Log(res.Created)
		t.Log(res.Choices[0].Delta)
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
