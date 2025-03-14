package main

import (
	"chatgpt-proxy/pkg/config"
	"chatgpt-proxy/routers"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func customRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if p := recover(); p != nil {
				if err, ok := p.(error); ok {
					// ignore panic abort handler for text/event-stream SSE
					if errors.Is(err, http.ErrAbortHandler) {
						return
					}
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// go run ./main.go --config=dev.config.yaml
func main() {
	cnf := config.GetConf()
	gin.SetMode(gin.ReleaseMode)
	// fmt.Printf("%+v\n", cnf)
	r := gin.Default()
	r.Use(customRecoveryMiddleware())
	routers.InitRouters(r)
	r.Run(fmt.Sprintf("%s:%d", cnf.Http.Host, cnf.Http.Port))
}
