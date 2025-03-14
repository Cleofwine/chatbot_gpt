package cron

import (
	"chatgpt-crontab/internal/wx/wecom"
	"chatgpt-crontab/internal/wx/wxofficial"
	"chatgpt-crontab/pkg/config"
	"chatgpt-crontab/pkg/log"
	"fmt"

	"github.com/robfig/cron/v3"
)

// func Run() {
// 	// cron.WithSeconds() 启用每秒定时任务
// 	// cron.WithLocation(time.Local) 设置时区为本地时区
// 	// c := cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
// 	c := cron.New(cron.WithSeconds())
// 	c.AddFunc("* * * * * *", func() {
// 		fmt.Println("每秒任务")
// 	})
// 	c.Run()
// }

func Run() {
	cnf := config.GetConf()
	c := cron.New()
	c.AddFunc("* * * * *", func() {
		fmt.Println("每分钟任务")
		for _, item := range cnf.WeComs {
			weCom := wecom.NewWecom(item.CorpId, item.CorpSecret, item.App)
			err := weCom.RefreshToken()
			if err != nil {
				log.Error(err)
				continue
			}
		}
		for _, item := range cnf.WxOfficials {
			official := wxofficial.NewWxOfficial(item.AppId, item.AppSecret)
			err := official.RefreshToken()
			if err != nil {
				log.Error(err)
				continue
			}
		}
	})
	c.Run()
}
