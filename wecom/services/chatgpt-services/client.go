package chatgptservices

import (
	"chatgpt-wecom/pkg/config"
	"chatgpt-wecom/services"
	"sync"

	defaultclient "chatgpt-wecom/services/default_client"
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
		pool = c.GetPool(cnf.DependOnServices.ChatgptService.Address)
	})
	return pool
}
