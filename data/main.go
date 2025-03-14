package main

import (
	"chatgpt-data/data"
	"chatgpt-data/pkg/cmd"
	"chatgpt-data/pkg/config"
	"chatgpt-data/pkg/db/mysql"
	"chatgpt-data/pkg/log"
	"chatgpt-data/proto"
	"chatgpt-data/server"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	loadDependOn()
	cnf := config.GetConf()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(server.GetOptions()...)
	logger := log.NewLogger()
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)
	chatRecordsData := data.NewChatRecordsData(mysql.GetDB(), logger)
	proto.RegisterChatGPTDataServer(s, server.NewChatgptDataServer(cnf, logger, chatRecordsData))

	// 使用同一个tcp连接，属于tcp的多路复用，因为使用了同一个server，这里也被拦截器鉴权拦截，需要修改拦截器
	healthCheck := health.NewServer()
	grpchealth.RegisterHealthServer(s, healthCheck)

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func loadDependOn() {
	config.InitConf(cmd.Args.Config)
	cnf := config.GetConf()
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetLevel(cnf.Log.Level)
	mysql.InitMysql()
}
