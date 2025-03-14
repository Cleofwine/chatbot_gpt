package datas

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/services"
	defaultclient "chatgpt-service/services/default_client"
	"sync"
)

var pool services.ClientPool
var once sync.Once

type chatGPTDataClient struct {
	defaultclient.DefaultClient
}

func GetChatGPTDataClientPool() services.ClientPool {
	once.Do(func() {
		cnf := config.GetConf()
		c := &chatGPTDataClient{}
		pool = c.GetPool(cnf.DependOnServices.ChatgptData.Address)
	})
	return pool
}
