package middleware

import (
	"chatgpt-proxy/pkg/config"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cnf := config.GetConf()
		authorization := ctx.Request.Header.Get("Authorization")
		confAuthorization := fmt.Sprintf("Bearer %s", cnf.Http.AccessToken)
		if authorization != confAuthorization {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(len(cnf.Chat.APIKeys))
		apiKey := cnf.Chat.APIKeys[randIndex]
		newAPIKEY := fmt.Sprintf("Bearer %s", apiKey)
		ctx.Request.Header.Set("Authorization", newAPIKEY)
		ctx.Next()
	}
}
