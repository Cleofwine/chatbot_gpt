package testcase

import (
	"chatgpt-data/data"
	"chatgpt-data/pkg/config"
	"chatgpt-data/pkg/db/mysql"
	"chatgpt-data/pkg/log"
	"chatgpt-data/proto"
	"chatgpt-data/server"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"net"
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

// go test ./test-case/. --config=../dev.config.yaml
func TestMain(m *testing.M) {
	flag.Parse()
	config.InitConf(*configPath)
	mysql.InitMysql()
	m.Run()
}

func TestAddRecord(t *testing.T) {
	dataList := []*proto.Record{
		{
			Account:         "1111",
			GroupId:         "AAA",
			UserMsg:         "你好，今天天气好吗",
			UserMsgTokens:   20,
			UserMsgKeywords: []string{"你好", "天气"},
			AiMsg:           "可以的，昨天也不错啊",
			AiMsgTokens:     50,
			ReqTokens:       100,
			CreateAt:        123456789,
			EndpointAccount: "4044785548",
			EnterpriseId:    "lin",
			Endpoint:        0,
		},
		{
			Account:         "22222",
			GroupId:         "BBBBBB",
			UserMsg:         "你好，今天天气好吗",
			UserMsgTokens:   20,
			UserMsgKeywords: []string{"你好", "天气"},
			AiMsg:           "可以的，昨天也不错啊",
			AiMsgTokens:     50,
			ReqTokens:       100,
			CreateAt:        123456789,
			EndpointAccount: "4044785999",
			EnterpriseId:    "lin",
			Endpoint:        1,
		},
		{
			Account:         "333333",
			GroupId:         "DDDDD",
			UserMsg:         "你好，今天天气好吗",
			UserMsgTokens:   20,
			UserMsgKeywords: []string{"你好", "天气"},
			AiMsg:           "可以的，昨天也不错啊",
			AiMsgTokens:     50,
			ReqTokens:       100,
			CreateAt:        123456789,
			EndpointAccount: "4044785777",
			EnterpriseId:    "linCC",
			Endpoint:        2,
		},
	}
	// server
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		t.Error(err)
		return
	}
	// 添加服务器启动选项
	s := grpc.NewServer(getServerOptions()...)
	defer s.Stop()
	conf := config.GetConf()
	logger := log.NewLogger()
	chatRecordData := data.NewChatRecordsData(mysql.GetDB(), logger)
	chatGPTDataServer := server.NewChatgptDataServer(conf, logger, chatRecordData)
	proto.RegisterChatGPTDataServer(s, chatGPTDataServer)
	go func() {
		err = s.Serve(lis)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	// client
	conn, err := grpc.Dial("localhost:50051", getClientOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	client := proto.NewChatGPTDataClient(conn)
	ctx := context.Background()
	// insecure鉴权方式下加入bearer token的方法
	ctx = appendBearerTokenToContext(ctx)
	for _, item := range dataList {
		res, err := client.AddRecord(ctx, item)
		if err != nil {
			t.Error(err)
			return
		}
		record, err := chatRecordData.GetById(res.GetId())
		if err != nil {
			t.Error(err)
			return
		}
		if record.CreateAt != item.CreateAt || record.ReqTokens != int(item.ReqTokens) ||
			record.AIMsg != item.AiMsg || record.AIMsgTokens != int(item.AiMsgTokens) ||
			record.UserMsg != item.UserMsg || len(record.UserMsgKeywords) != len(item.UserMsgKeywords) ||
			record.UserMsgTokens != int(item.UserMsgTokens) || record.Account != item.Account ||
			record.GroupID != item.GroupId || record.EnterpriseId != item.EnterpriseId ||
			record.EndpointAccount != item.EndpointAccount || record.Endpoint != int(item.Endpoint) {
			t.Error("写入的记录与读取的不匹配")
		}
	}
}

func getClientOptions() []grpc.DialOption {
	clientOption := make([]grpc.DialOption, 0)
	clientOption = append(clientOption, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// clientOption = append(clientOption, getTlsOpt("../server/cert_ex/ca_cert.pem", "chatgpt-data.grpc.lin.com"))
	// clientOption = append(clientOption, getMTLSOpt("../server/cert_ex/ca_cert.pem", "../server/cert_ex/client_cert.pem", "../server/cert_ex/client_key.pem"))
	// clientOption = append(clientOption, getAuth())
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
	// opts = append(opts, server.GetTlsOpt("../server/cert_ex/server_cert.pem", "../server/cert_ex/server_key.pem"))
	// opts = append(opts, server.GetMTLSOpt("../server/cert_ex/client_ca_cert.pem", "../server/cert_ex/server_cert.pem", "../server/cert_ex/server_key.pem"))
	opts = append(opts, grpc.StreamInterceptor(server.StreamInterceptor))
	opts = append(opts, grpc.UnaryInterceptor(server.UnaryInterceptor))
	return opts
}
