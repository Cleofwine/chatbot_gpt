package chatgptservices

import (
	"chatgpt-qq/config"
	"chatgpt-qq/services"
	defaultclient "chatgpt-qq/services/default_client"
	"sync"
)

var pool services.ClientPool
var once sync.Once

type chatGPTServicesClient struct {
	defaultclient.DefaultClient
}

func GetChatGPTServicesClientPool() services.ClientPool {
	once.Do(func() {
		cnf := config.GetConf()
		c := &chatGPTServicesClient{}
		pool = c.GetPool(cnf.ChatGPTService.Address)
	})
	return pool
}
