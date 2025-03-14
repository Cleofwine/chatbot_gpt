package chatgptcrontab

import (
	"sync"
	"wxofficial/pkg/config"
	"wxofficial/services"
	defaultclient "wxofficial/services/default_client"
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
