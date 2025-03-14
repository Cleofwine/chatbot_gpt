package server

import (
	"chatgpt-crontab/internal/wx"
	"chatgpt-crontab/internal/wx/wecom"
	"chatgpt-crontab/internal/wx/wxofficial"
	"chatgpt-crontab/pkg/config"
	"chatgpt-crontab/pkg/log"
	"chatgpt-crontab/proto"
	"context"
)

type tokenServer struct {
	proto.UnimplementedTokenServer
	config *config.Config
	log    *log.Logger
}

func NewTokenServer(config *config.Config, log *log.Logger) proto.TokenServer {
	return &tokenServer{
		config: config,
		log:    log,
	}
}

func (s *tokenServer) GetToken(ctx context.Context, in *proto.TokenRequest) (*proto.TokenResponse, error) {
	var token wx.Token
	secret := s.getSecret(in)
	if in.Typ == proto.TokenType_WECHATOFFICIAL {
		token = wxofficial.NewWxOfficial(in.Id, secret)
	} else if in.Typ == proto.TokenType_WECOM {
		token = wecom.NewWecom(in.Id, secret, in.App)
	}
	if token != nil {
		accessToken, err := token.GetToken()
		if err != nil {
			s.log.Error(err)
			return nil, err
		}
		res := &proto.TokenResponse{
			AccessToken: accessToken.AccessToken,
		}
		return res, err
	}
	return nil, nil
}

func (s *tokenServer) getSecret(in *proto.TokenRequest) string {
	if in.Typ == proto.TokenType_WECHATOFFICIAL {
		for _, item := range s.config.WxOfficials {
			if item.AppId == in.Id {
				return item.AppSecret
			}
		}
	} else if in.Typ == proto.TokenType_WECOM {
		for _, item := range s.config.WeComs {
			if item.CorpId == in.Id && item.App == in.App {
				return item.CorpSecret
			}
		}
	}
	return ""
}
