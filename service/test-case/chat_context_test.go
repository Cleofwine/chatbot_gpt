package testcase

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/proto"
	chatcontext "chatgpt-service/server/chat-context"
	"testing"

	"github.com/sashabaranov/go-openai"
)

// go test ./test-case/. -run ^TestWebContext$ --config=../dev.config.yaml -v
func TestWebContext(t *testing.T) {
	dataList := []*chatcontext.ChatMessage{
		{
			ID: "121314",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
			PID: "",
		},
		{
			ID: "222222",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你的吗？",
			},
			PID: "121314",
		},
		{
			ID: "33333",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
			PID: "222222",
		},
		{
			ID: "4444444",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你的吗？",
			},
			PID: "33333",
		},
		{
			ID: "5555555",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
			PID: "4444444",
		},
		{
			ID: "66666666",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你的吗？",
			},
			PID: "5555555",
		},
	}
	cache := chatcontext.GetCacheContext(proto.ChatEndpoint_WEB)
	for _, item := range dataList {
		err := cache.Set(item.ID, "", proto.ChatEndpoint_WEB, item, config.GetConf().Chat.ContextTTL)
		if err != nil {
			t.Error(err)
			return
		}
	}
	l, err := cache.Get("66666666", "", proto.ChatEndpoint_WEB)
	if err != nil {
		t.Error(err)
		return
	}
	cmList, ok := l.([]*chatcontext.ChatMessage)
	if !ok {
		t.Error("类型不对")
	}
	if len(cmList) == 0 || len(cmList) > config.GetConf().Chat.ContextLen {
		t.Error("上下文条目不对")
		t.Log(len(cmList))
	}
	for _, item := range cmList {
		t.Log(*item)
	}
	// t.Log(cmList)
}

// go test ./test-case/. -run ^TestQQContext$ --config=../dev.config.yaml -v
func TestQQContext(t *testing.T) {
	dataList := []*chatcontext.ChatMessage{
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好6",
			},
		},
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你的吗？5",
			},
		},
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好4",
			},
		},
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你的吗？3",
			},
		},
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好2",
			},
		},
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你的吗？1",
			},
		},
	}
	id := "447895125"
	group := "geirnpoin"
	cache := chatcontext.GetCacheContext(proto.ChatEndpoint_QQ)
	err := cache.Set(id, group, proto.ChatEndpoint_QQ, dataList, config.GetConf().Chat.ContextTTL)
	if err != nil {
		t.Error(err)
		return
	}
	l, err := cache.Get(id, group, proto.ChatEndpoint_QQ)
	if err != nil {
		t.Error(err)
		return
	}
	cmList, ok := l.([]*chatcontext.ChatMessage)
	if !ok {
		t.Error("类型不对")
		return
	}
	if len(cmList) == 0 || len(cmList) > config.GetConf().Chat.ContextLen {
		t.Error("上下文条目不对")
		t.Log(len(cmList))
		return
	}
	for _, item := range cmList {
		t.Log(*item)
	}
}
