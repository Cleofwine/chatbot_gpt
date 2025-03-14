package chatcontext

import (
	"chatgpt-service/pkg/db/redis"
	"chatgpt-service/proto"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

const KEYPREFIX = "context_"

func GetKey(id, groupid string, endpoint proto.ChatEndpoint) string {
	return redis.GetKey(fmt.Sprintf("%s_%s_%s_%d", KEYPREFIX, id, groupid, endpoint))
}

type ChatMessage struct {
	// 当前记录的ID
	ID string `json:"id,omitempty"`
	// 上一条记录的ID
	PID string `json:"pid,omitempty"`
	// 消息内容
	Message openai.ChatCompletionMessage `json:"msg,omitempty"`
	// 该消息token数
	TokensNum int `json:"t_num,omitempty"`
}

type ICacheContext interface {
	Get(id, group string, endpoint proto.ChatEndpoint) (interface{}, error)
	Set(id, group string, endpoint proto.ChatEndpoint, value interface{}, ttl int) error
}

func GetCacheContext(endpoint proto.ChatEndpoint) ICacheContext {
	if endpoint == proto.ChatEndpoint_WEB {
		return &notHasID{}
	} else {
		return &hasID{}
	}
}
