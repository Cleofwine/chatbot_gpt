package main

import (
	"chatgpt-crontab/cron"
	"chatgpt-crontab/pkg/cmd"
	"chatgpt-crontab/pkg/config"
	"chatgpt-crontab/pkg/db/redis"
	"chatgpt-crontab/pkg/log"
	"chatgpt-crontab/proto"
	"chatgpt-crontab/token-server/server"
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
		log.Error(err)
		panic(err)
	}
	s := grpc.NewServer(server.GetOptions()...)

	logger := log.NewLogger()
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)
	tokenServer := server.NewTokenServer(cnf, logger)
	proto.RegisterTokenServer(s, tokenServer)

	// 添加健康检查逻辑
	healthCheck := health.NewServer()
	grpchealth.RegisterHealthServer(s, healthCheck)

	// 定时任务，刷新token
	go cron.Run()

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func loadDependOn() {
	// 初始化配置
	config.InitConf(cmd.Args.Config)
	cnf := config.GetConf()
	// fmt.Printf("%+v\n", cnf)

	// 初始化日志
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetLevel(cnf.Log.Level)

	// 初始化redis
	redis.InitRedisPool()
}
