package main

import (
	"chatgpt-sensitive/pkg/cmd"
	"chatgpt-sensitive/pkg/config"
	"chatgpt-sensitive/pkg/filter"
	"chatgpt-sensitive/proto"
	"chatgpt-sensitive/server"
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
	proto.RegisterChatGPTSensitiveServer(s, server.NewSensitiveWordServer(filter.GetFilter()))

	// 使用同一个tcp连接，属于tcp的多路复用，因为使用了同一个server，这里也被拦截器鉴权拦截，需要修改拦截器
	healthCheck := health.NewServer()
	grpchealth.RegisterHealthServer(s, healthCheck)

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func loadDependOn() {
	config.InitConf(cmd.Args.Config)
	filter.InitFilter(cmd.Args.Dict)
}
