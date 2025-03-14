package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 如何找到适合的rate和b？
// 我们进行了压测（bench测试），1000x请求，cpu1000是撑不住的(一个是代理代码处崩溃，一个是gpt返回429说明被官方限流)，而100x请求，cpu100可以撑住（均返回200）。那么2核的机器，每核每s进行50次，我们考虑容器部署，每个容器分配0.3个核，那么0.3核1s分配18次左右的请求，又考虑到我们的业务io输出可能时间较长，那么我们设置初始速率为10，桶大小也为10。也就是每秒最多并发10个请求，下一秒就会恢复满。
// 如何找到更精确的值？继续进行压测，找到极限的值，多一点点就会撑不住，少一点点刚好不出问题。
func RateLimit(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)
	return func(ctx *gin.Context) {
		if !limiter.Allow() {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
