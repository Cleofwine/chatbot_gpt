package official

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
	"wxofficial/pkg/config"
	"wxofficial/pkg/log"
	chatgptservices "wxofficial/services/chatgpt-services"
	chatgpt_service_proto "wxofficial/services/chatgpt-services/proto"
	crontab "wxofficial/services/crontab"
	crontab_proto "wxofficial/services/crontab/proto"
	defaultclient "wxofficial/services/default_client"

	"github.com/gin-gonic/gin"
)

type Message struct {
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   int64
	MsgType      CDATA
	Content      CDATA
	MsgId        int64
	MsgDataId    int64
	Idx          int64
	// 图片
	PicUrl  CDATA
	MediaId CDATA
	// 语音消息
	Format      CDATA
	Recognition CDATA
}

type CDATA struct {
	Value string `xml:",cdata"`
}

type AutoReplyTextMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   int64
	MsgType      CDATA
	Content      CDATA
}

func CheckSignature(ctx *gin.Context) {
	cnf := config.GetConf()
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	sign := makeSignature(cnf.Official.Token, timestamp, nonce)
	echoStr := ctx.Query("echostr")
	if signature == sign {
		// 回复
		ctx.Data(200, "text/plain;charset=utf-8;", []byte(echoStr))
		return
	}
}

func ReceiveMessage(ctx *gin.Context) {
	cnf := config.GetConf()
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	sign := makeSignature(cnf.Official.Token, timestamp, nonce)
	if signature != sign {
		ctx.Data(200, "text/plain;charset=utf-8;", []byte("success"))
		return
	}
	// 获取body内容
	msg := &Message{}
	err := ctx.BindXML(msg)
	if err != nil {
		log.Error(err)
		ctx.Data(200, "text/plain;charset=utf-8", []byte("success"))
		return
	}
	if msg.MsgType.Value != "text" {
		replyMsg := AutoReplyTextMessage{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      CDATA{Value: "text"},
			Content:      CDATA{Value: "抱歉，目前不支持文本消息以外的其他消息类型"},
		}
		ctx.XML(200, replyMsg)
		return
	}
	ctx.Data(200, "text/plain;charset=utf-8", []byte("success"))
	// 接入chatgpt
	go func() {
		content, err := generateChatCompletion(msg.FromUserName.Value, msg.ToUserName.Value, msg.Content.Value)
		if err != nil {
			log.Error(err)
			return
		}
		if content == "" {
			return
		}
		sendKfTextMsg(msg.FromUserName.Value, content)
	}()
}

func makeSignature(token, timestamp, nonce string) string {
	sortArr := []string{
		token, timestamp, nonce,
	}
	sort.Strings(sortArr)
	var buffer bytes.Buffer
	for _, value := range sortArr {
		buffer.WriteString(value)
	}
	sha := sha1.New()
	sha.Write(buffer.Bytes())
	signature := fmt.Sprintf("%x", sha.Sum(nil))
	return signature
}

type KfTextMessage struct {
	ToUser  string `json:"touser"`
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

type SendResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func sendKfTextMsg(toUser, content string) {
	// 获取accessToken
	accessToken := getAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + accessToken
	method := "post"
	replyMsg := &KfTextMessage{
		ToUser:  toUser,
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{Content: content},
	}
	payloadBytes, err := json.Marshal(replyMsg)
	if err != nil {
		log.Error(err)
		return
	}
	payload := strings.NewReader(string(payloadBytes))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return
	}
	sendRes := &SendResponse{}
	err = json.Unmarshal(body, sendRes)
	if err != nil {
		log.Error(err)
		return
	}
	if sendRes.ErrCode != 0 {
		log.Error(sendRes.ErrMsg)
		return
	}
}

func generateChatCompletion(userID, endpointAccount, content string) (string, error) {
	cnf := config.GetConf()
	chatGPTServiceClientPool := chatgptservices.GetChatGPTServicesClientPool()
	conn := chatGPTServiceClientPool.Get()
	defer chatGPTServiceClientPool.Put(conn)

	client := chatgpt_service_proto.NewChatGPTServiceServerClient(conn)
	ctx1 := context.Background()
	ctx1 = defaultclient.AppendBearerTokenToContext(ctx1, cnf.DependOnServices.ChatgptService.AccessToken)
	in := &chatgpt_service_proto.ChatCompletionReq{
		Id:              userID,
		Message:         content,
		Endpoint:        chatgpt_service_proto.ChatEndpoint_WECHATOFFICAL,
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

func getAccessToken() string {
	cnf := config.GetConf()
	crontabClientPool := crontab.GetChatGPTCrontabClientPool()
	conn := crontabClientPool.Get()
	defer crontabClientPool.Put(conn)

	client := crontab_proto.NewTokenClient(conn)
	in := &crontab_proto.TokenRequest{
		Typ: crontab_proto.TokenType_WECHATOFFICIAL,
		Id:  cnf.Official.AppId,
		App: "",
	}
	ctx := context.Background()
	ctx = defaultclient.AppendBearerTokenToContext(ctx, cnf.DependOnServices.ChatgptCrontab.AccessToken)
	res, err := client.GetToken(ctx, in)
	if err != nil {
		log.Error(err)
		return ""
	}
	return res.AccessToken
}
