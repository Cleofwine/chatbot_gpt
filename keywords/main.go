package main

import (
	"chatgpt-keywords/pkg/cmd"
	"chatgpt-keywords/pkg/config"
	"chatgpt-keywords/pkg/filter"
	"chatgpt-keywords/proto"
	"chatgpt-keywords/server"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
)

// 初始化词库
//  go run . --config=dev.config.yaml --dict=cainiao-coding.txt --init-dict=true

func main() {
	loadDependOn()
	if cmd.Args.InitDict {
		filter.OverwriteDict(cmd.Args.Dict)
		return
	}

	cnf := config.GetConf()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(server.GetOptions()...)
	proto.RegisterChatGPTKeywordsServer(s, server.NewKeyWordsServer(filter.GetFilter()))

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
