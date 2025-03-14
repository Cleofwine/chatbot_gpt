package routers

import (
	"chatgpt-proxy/health"
	"chatgpt-proxy/middleware"
	"chatgpt-proxy/proxy"

	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	r.GET("/health", health.Health)
	r.Use(middleware.Auth(), middleware.RateLimit(10, 10))
	initProxyRouter(r)
	// middleCaseRouter(r)
}

func initProxyRouter(r *gin.Engine) {
	p := proxy.NewProxy()
	v1 := r.Group("v1")
	v1.Any("/*relativePath", p.ChatProxy)
}

/* 下面是一个添加中间件的案例 */
// func middleCaseRouter(r *gin.Engine) {
// 	v2 := r.Group("/v2")
// 	v2.Use(func(ctx *gin.Context) { // 在这个位置添加意味着v2这个对象后面的所有接口都应用该中间件
// 		fmt.Println("call middleware 1")
// 	})
// 	v2.GET("/view", view)

// 	v3 := r.Group("/v2") // 路由可以是一样的，但是对应的对象不一样，就可以应用不同的中间件，如不同的鉴权方式
// 	v3.GET("/personal", func(ctx *gin.Context) {
// 		fmt.Println("call middleware 2")
// 	}, personal) // 这个位置添加中间件意味着，v3这个对象中只有这一个路由对应的位置添加该中间件
// 	v3.GET("/personal1", personal)
// }
// func view(ctx *gin.Context) {
// 	fmt.Println("call view")
// }
// func personal(ctx *gin.Context) {
// 	fmt.Println("call personal")
// }
