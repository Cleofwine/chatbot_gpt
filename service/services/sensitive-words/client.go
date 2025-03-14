package sensitive

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/services"
	defaultclient "chatgpt-service/services/default_client"
	"sync"
)

var pool services.ClientPool
var once sync.Once

type sensitiveClient struct {
	defaultclient.DefaultClient
}

func GetSensitiveClientPool() services.ClientPool {
	once.Do(func() {
		cnf := config.GetConf()
		c := &sensitiveClient{}
		pool = c.GetPool(cnf.DependOnServices.ChatgptSensitive.Address)
	})
	return pool
}
