package server

import (
	"chatgpt-sensitive/pkg/filter"
	"chatgpt-sensitive/proto"
	"context"
)

type SensitiveWordServer struct {
	proto.UnimplementedChatGPTSensitiveServer
	filter filter.ISensitiveFilter
}

func NewSensitiveWordServer(filter filter.ISensitiveFilter) proto.ChatGPTSensitiveServer {
	return &SensitiveWordServer{
		filter: filter,
	}
}

func (s SensitiveWordServer) Validate(ctx context.Context, in *proto.ValidateReq) (*proto.ValidateRes, error) {
	ok, word := s.filter.Validate(in.Text)
	return &proto.ValidateRes{
		Ok:   ok,
		Word: word,
	}, nil
}
