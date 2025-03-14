package server

import (
	"chatgpt-data/data"
	"chatgpt-data/pkg/config"
	"chatgpt-data/pkg/log"
	"chatgpt-data/proto"
	"context"
	"time"
)

type ChatgptDataServer struct {
	proto.UnimplementedChatGPTDataServer
	// 配置信息
	config *config.Config
	// log信息
	log *log.Logger
	// 数据访问层
	chatRecordsData data.IChatRecordsData
}

func NewChatgptDataServer(conf *config.Config, log *log.Logger, chatRecordsData data.IChatRecordsData) *ChatgptDataServer {
	return &ChatgptDataServer{
		config:          conf,
		log:             log,
		chatRecordsData: chatRecordsData,
	}
}

func (s *ChatgptDataServer) AddRecord(ctx context.Context, in *proto.Record) (*proto.RecordRes, error) {
	cr := &data.ChatRecord{}
	cr.Account = in.Account
	cr.GroupID = in.GroupId
	cr.UserMsg = in.UserMsg
	cr.UserMsgTokens = int(in.UserMsgTokens)
	cr.AIMsg = in.AiMsg
	cr.AIMsgTokens = int(in.AiMsgTokens)
	cr.UserMsgKeywords = in.UserMsgKeywords
	cr.CreateAt = in.CreateAt
	cr.ReqTokens = int(in.ReqTokens)
	cr.EnterpriseId = in.EnterpriseId
	cr.Endpoint = int(in.Endpoint)
	cr.EndpointAccount = in.EndpointAccount
	if cr.CreateAt == 0 {
		cr.CreateAt = time.Now().Unix()
	}
	err := s.chatRecordsData.Add(cr)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	out := &proto.RecordRes{
		Id: cr.ID,
	}
	return out, err
}
