package chatgptservices

import (
	"chatgpt-web/pkg/config"
	"chatgpt-web/services"
	defaultclient "chatgpt-web/services/default_client"
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
		pool = c.GetPool(cnf.DependOnServices.ChatgptServices.Address)
	})
	return pool
}
