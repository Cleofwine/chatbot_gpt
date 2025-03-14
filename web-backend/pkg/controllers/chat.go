package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"chatgpt-web/pkg/config"
	"chatgpt-web/pkg/log"
	chatgptservices "chatgpt-web/services/chatgpt-services"
	chatgptservices_proto "chatgpt-web/services/chatgpt-services/proto"
	defaultclient "chatgpt-web/services/default_client"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ChatPrimedTokens = 2
)

type ChatService struct {
	config *config.Config
	log    *log.Logger
}

type ChatMessageRequest struct {
	Prompt  string                    `json:"prompt"`
	Options ChatMessageRequestOptions `json:"options"`
}

type ChatMessageRequestOptions struct {
	Name            string `json:"name"`
	ParentMessageId string `json:"parentMessageId"`
}

type ChatMessage struct {
	ID              string                                         `json:"id"`
	Text            string                                         `json:"text"`
	Role            string                                         `json:"role"`
	Name            string                                         `json:"name"`
	Delta           string                                         `json:"delta"`
	Detail          *chatgptservices_proto.ChatCompletionStreamRes `json:"detail"`
	TokenCount      int                                            `json:"tokenCount"`
	ParentMessageId string                                         `json:"parentMessageId"`
}

func NewChatService(config *config.Config, log *log.Logger) (*ChatService, error) {
	chat := ChatService{
		config: config,
		log:    log,
	}
	return &chat, nil
}

const ChatMessageRoleAssistant = "assistant"

func (chat *ChatService) ChatProcess(ctx *gin.Context) {
	payload := ChatMessageRequest{}
	if err := ctx.BindJSON(&payload); err != nil {
		chat.log.Error(err)
		ctx.JSON(200, gin.H{
			"status":  "Fail",
			"message": fmt.Sprintf("%v", err),
			"data":    nil,
		})
		return
	}
	// username := ctx.GetString("username")

	messageID := uuid.New().String()

	result := ChatMessage{
		ID:              uuid.New().String(),
		Role:            ChatMessageRoleAssistant,
		Text:            "",
		ParentMessageId: messageID,
	}

	// 接入核心服务

	chatGPTServicesClientPool := chatgptservices.GetChatGPTServicesClientPool()
	conn := chatGPTServicesClientPool.Get()
	defer chatGPTServicesClientPool.Put(conn)

	client := chatgptservices_proto.NewChatGPTServiceServerClient(conn)
	ctx1 := context.Background()
	ctx1 = defaultclient.AppendBearerTokenToContext(ctx1, chat.config.DependOnServices.ChatgptServices.AccessToken)
	in := &chatgptservices_proto.ChatCompletionReq{
		Id:              messageID,
		Message:         payload.Prompt,
		Pid:             payload.Options.ParentMessageId,
		Endpoint:        chatgptservices_proto.ChatEndpoint_WEB,
		EnterpriseId:    chat.config.Enterprise.Id,
		EnableContext:   false,
		EndpointAccount: chat.config.Enterprise.Id,
		ChatParam: &chatgptservices_proto.ChatParam{
			Model:             chat.config.Chat.Model,
			BotDesc:           chat.config.Chat.BotDesc,
			ContextLen:        int32(chat.config.Chat.ContextLen),
			MinResponseTokens: int32(chat.config.Chat.MinResponseTokens),
			ContextTTL:        int32(chat.config.Chat.ContextTTL),
			Temperature:       chat.config.Chat.Temperature,
			PresencePenalty:   chat.config.Chat.PresencePenalty,
			FrequencyPenalty:  chat.config.Chat.FrequencyPenalty,
			TopP:              chat.config.Chat.TopP,
			MaxTokens:         int32(chat.config.Chat.MaxTokens),
		},
	}
	if payload.Options.ParentMessageId != "" {
		in.EnableContext = true
	}
	stream, err := client.ChatCompletionStream(ctx1, in)
	if err != nil {
		chat.log.Error(err)
		ctx.JSON(200, gin.H{
			"status":  "Fail",
			"message": fmt.Sprintf("%v", err),
			"data":    nil,
		})
		return
	}
	defer stream.CloseSend()

	firstChunk := true

	ctx.Header("Content-type", "application/octet-stream")
	for {
		rsp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return
		}

		if err != nil {
			chat.log.Error(err)
			ctx.JSON(200, gin.H{
				"status":  "Fail",
				"message": fmt.Sprintf("OpenAI Event Error %v", err),
				"data":    nil,
			})
			return
		}

		if rsp.Id != "" {
			result.ID = rsp.Id
		}

		if len(rsp.Choices) > 0 {
			content := rsp.Choices[0].Delta.Content
			result.Delta = content
			if len(content) > 0 {
				result.Text += content
			}
			result.Detail = rsp
		}

		bts, err := json.Marshal(result)
		if err != nil {
			chat.log.Error(err)
			ctx.JSON(200, gin.H{
				"status":  "Fail",
				"message": fmt.Sprintf("OpenAI Event Marshal Error %v", err),
				"data":    nil,
			})
			return
		}

		if !firstChunk {
			ctx.Writer.Write([]byte("\n"))
		} else {
			firstChunk = false
		}

		if _, err := ctx.Writer.Write(bts); err != nil {
			chat.log.Error(err)
			return
		}

		ctx.Writer.Flush()
	}
}
