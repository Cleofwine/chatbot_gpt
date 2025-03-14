package chatgptcrontab

import (
	"chatgpt-wecom/pkg/config"
	"chatgpt-wecom/services"
	"sync"

	defaultclient "chatgpt-wecom/services/default_client"
)

var pool services.ClientPool
var once sync.Once

type chatGPTCrontabClient struct {
	defaultclient.DefaultClient
}

func GetChatGPTCrontabClientPool() services.ClientPool {
	once.Do(func() {
		cnf := config.GetConf()
		c := &chatGPTCrontabClient{}
		pool = c.GetPool(cnf.DependOnServices.ChatgptCrontab.Address)
	})
	return pool
}
