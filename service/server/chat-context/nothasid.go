package chatcontext

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/pkg/db/redis"
	"chatgpt-service/pkg/log"
	"chatgpt-service/proto"
	"context"
	"encoding/json"
	"errors"
	"time"

	redis2 "github.com/redis/go-redis/v9"
)

// web
// 没有登录信息，只能使用id和pid来链式查询上下文
type notHasID struct {
}

func (c *notHasID) Get(pid, group string, endpoint proto.ChatEndpoint) (value interface{}, err error) {
	maxLen := config.GetConf().Chat.ContextLen
	list := make([]*ChatMessage, 0)
	value = list
	var item *ChatMessage
	var i = 0
	for {
		key := GetKey(pid, group, endpoint)
		// fmt.Printf("get key = %s\n", key)
		item, err = c.getOnce(key)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if item == nil {
			// fmt.Printf("item == nil")
			break
		}
		list = append(list, item)
		pid = item.PID
		i++
		if i >= maxLen || pid == "" {
			break
		}
	}
	value = list
	return
}
func (c *notHasID) Set(id, group string, endpoint proto.ChatEndpoint, value interface{}, ttl int) error {
	key := GetKey(id, group, endpoint)
	item, ok := value.(*ChatMessage)
	if !ok {
		err := errors.New("value 类型不匹配")
		log.Error("value 类型不匹配")
		return err
	}
	redisPool := redis.GetPool()
	redisClient := redisPool.Get()
	defer redisPool.Put(redisClient)
	bytes, err := json.Marshal(item)
	if err != nil {
		log.Error(err)
		return err
	}
	err = redisClient.SetEx(context.Background(), key, string(bytes), time.Duration(ttl)*time.Second).Err()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c *notHasID) getOnce(key string) (msg *ChatMessage, err error) {
	redisPool := redis.GetPool()
	redisClient := redisPool.Get()
	defer redisPool.Put(redisClient)
	str, err := redisClient.Get(context.Background(), key).Result()
	if err == redis2.Nil {
		err = nil
		return
	}
	if err != nil {
		log.Error(err)
		return
	}
	value := &ChatMessage{}
	err = json.Unmarshal([]byte(str), value)
	if err != nil {
		log.Error(err)
		return
	}
	msg = value
	return
}
