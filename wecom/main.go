package main

import (
	"chatgpt-wecom/pkg/cmd"
	"chatgpt-wecom/pkg/config"
	"chatgpt-wecom/pkg/log"
	"chatgpt-wecom/routers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	loadDependOn()
	gin.New()
	r := gin.Default()
	routers.InitRouters(r)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	cnf := config.GetConf()
	addr := fmt.Sprintf("%s:%d", cnf.Http.Host, cnf.Http.Port)
	r.Run(addr)
}

func loadDependOn() {
	config.InitConf(cmd.Args.Config)
	cnf := config.GetConf()
	// fmt.Printf("%+v", cnf)
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetLevel(cnf.Log.Level)
}
