package msghandler

import (
	"chatgpt-wechat/pkg/config"
	"chatgpt-wechat/pkg/log"
	chatgptservices "chatgpt-wechat/services/chatgpt-services"
	chatgpt_service_proto "chatgpt-wechat/services/chatgpt-services/proto"
	defaultclient "chatgpt-wechat/services/default_client"
	"context"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

type msgHandler struct {
	// bot openwechat.Bot
}

func NewMsgHandler() *msgHandler {
	return &msgHandler{}
}

func (mh *msgHandler) TextHandler(ctx *openwechat.MessageContext) {
	var user *openwechat.User
	var group *openwechat.Group
	var content string
	var err error
	var groupID string
	content = ctx.Content
	user, err = ctx.Sender()
	if err != nil {
		return
	}
	// 判断是否是群聊
	if ctx.IsSendByGroup() {
		group = &openwechat.Group{User: user}
		if ctx.IsAt() {
			return
		}

		user, err = ctx.SenderInGroup()
		if err != nil {
			log.Error(err)
			return
		}
		content = strings.TrimSpace(strings.ReplaceAll(content, "@"+user.Self().NickName, ""))
		groupID = group.ID()
	}
	// fmt.Println("当前请求的用户信息：", user.ID(), user.UserName, user.NickName)
	// fmt.Println("接受请求的用户信息：", user.Self().ID(), user.Self().UserName, user.Self().NickName)
	// if group != nil {
	// 	fmt.Println("发送请求的群组信息：", group.ID(), group.UserName, group.NickName)
	// }
	// 从chatgpt获取回复
	replyText, err := generateChatCompletion(user.UserName, groupID, string(user.Self().ID()), content)
	if err != nil {
		log.Error(err)
	}
	if ctx.IsSendByGroup() {
		replyText = "@" + user.NickName + " " + replyText
	}
	ctx.ReplyText(replyText)
}

func generateChatCompletion(userID, groupID, endpointAccount, content string) (string, error) {
	cnf := config.GetConf()
	chatGPTServiceClientPool := chatgptservices.GetChatGPTServicesClientPool()
	conn := chatGPTServiceClientPool.Get()
	defer chatGPTServiceClientPool.Put(conn)

	client := chatgpt_service_proto.NewChatGPTServiceServerClient(conn)
	ctx1 := context.Background()
	ctx1 = defaultclient.AppendBearerTokenToContext(ctx1, cnf.DependOnServices.ChatgptService.AccessToken)
	in := &chatgpt_service_proto.ChatCompletionReq{
		Id:              userID,
		GroupId:         groupID,
		Message:         content,
		Endpoint:        chatgpt_service_proto.ChatEndpoint_WECHAT,
		EnterpriseId:    cnf.Enterprise.Id,
		EnableContext:   cnf.Chat.EnableContext,
		EndpointAccount: endpointAccount,
		ChatParam: &chatgpt_service_proto.ChatParam{
			Model:             cnf.Chat.Model,
			BotDesc:           cnf.Chat.BotDesc,
			ContextLen:        int32(cnf.Chat.ContextLen),
			MinResponseTokens: int32(cnf.Chat.MinResponseTokens),
			ContextTTL:        int32(cnf.Chat.ContextTTL),
			Temperature:       cnf.Chat.Temperature,
			PresencePenalty:   cnf.Chat.PresencePenalty,
			FrequencyPenalty:  cnf.Chat.FrequencyPenalty,
			TopP:              cnf.Chat.TopP,
			MaxTokens:         int32(cnf.Chat.MaxTokens),
		},
	}
	res, err := client.ChatCompletion(ctx1, in)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return res.Choices[0].Message.Content, nil
}
