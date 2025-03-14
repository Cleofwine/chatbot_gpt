package server

import (
	"chatgpt-keywords/pkg/filter"
	"chatgpt-keywords/proto"
	"context"
)

type keywordsServer struct {
	proto.UnimplementedChatGPTKeywordsServer
	filter filter.IKeyworsFilter
}

func NewKeyWordsServer(filter filter.IKeyworsFilter) proto.ChatGPTKeywordsServer {
	return &keywordsServer{
		filter: filter,
	}
}

func (s keywordsServer) FindAll(ctx context.Context, in *proto.FindAllReq) (*proto.FindAllRes, error) {
	list := s.filter.FindAll(in.Text)
	return &proto.FindAllRes{
		Words: list,
	}, nil
}
