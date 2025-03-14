package redis

import (
	"chatgpt-service/pkg/config"
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type redisPool struct {
	pool sync.Pool
}

var pool RedisPool

type RedisPool interface {
	Get() *redis.Client
	Put(client *redis.Client)
}

func (p *redisPool) Get() *redis.Client {
	client := p.pool.Get().(*redis.Client)
	if client.Ping(context.Background()).Err() != nil {
		client = p.pool.New().(*redis.Client)
	}
	return client
}
func (p *redisPool) Put(client *redis.Client) {
	if client.Ping(context.Background()).Err() != nil {
		return
	}
	p.pool.Put(client)
}

func InitRedisPool() {
	pool = getPool()

}

func getPool() RedisPool {
	return &redisPool{
		pool: sync.Pool{
			New: func() any {
				conf := config.GetConf()
				rdb := redis.NewClient(&redis.Options{
					Addr:     fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
					Password: conf.Redis.PWD,
				})
				return rdb
			},
		},
	}
}

func GetPool() RedisPool {
	return pool
}
