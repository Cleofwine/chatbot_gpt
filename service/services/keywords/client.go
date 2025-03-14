package keywords

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/services"
	defaultclient "chatgpt-service/services/default_client"
	"sync"
)

var pool services.ClientPool
var once sync.Once

type keywordsClient struct {
	defaultclient.DefaultClient
}

func GetKeywordsClientPool() services.ClientPool {
	once.Do(func() {
		cnf := config.GetConf()
		c := &keywordsClient{}
		pool = c.GetPool(cnf.DependOnServices.ChatgptKeywords.Address)
	})
	return pool
}
