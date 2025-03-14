package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"chatgpt-web/pkg/cmd"
	"chatgpt-web/pkg/config"
	"chatgpt-web/pkg/controllers"
	"chatgpt-web/pkg/log"
	"chatgpt-web/pkg/middlewares"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

type ChatGPTWebServer struct {
	config *config.Config
	log    *log.Logger
}

func NewChatGPTWebServer(config *config.Config, log *log.Logger) *ChatGPTWebServer {
	return &ChatGPTWebServer{
		config: config,
		log:    log,
	}
}

func (r *ChatGPTWebServer) Run(ctx context.Context) error {
	gin.SetMode(gin.DebugMode)
	// if err := r.updateAssetsFiles(); err != nil {
	// 	return err
	// }
	go r.httpServer(ctx)
	return nil
}

func (r *ChatGPTWebServer) httpServer(ctx context.Context) {
	// accountService, err := controllers.NewAccountService(r.DataBase, r.BasicAuthUser, r.BasicAuthPassword)
	// if err != nil {
	// 	klog.Fatal(err)
	// }
	chatService, err := controllers.NewChatService(r.config, r.log)
	if err != nil {
		klog.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%d", r.config.Http.Host, r.config.Http.Port)
	r.log.InfoF("ChatGPT Web Server on: %s", addr)
	fmt.Printf("ChatGPT Web Server on: %s", addr)
	server := &http.Server{
		Addr: addr,
	}
	// entry, proxy := gin.Default(), gin.Default()
	entry := gin.Default()
	chat := entry.Group("/api")
	entry.Use(middlewares.Cors())
	if len(r.config.Http.BasicAuthUser) > 0 {
		accounts := gin.Accounts{}
		users := strings.Split(r.config.Http.BasicAuthUser, ",")
		passwords := strings.Split(r.config.Http.BasicAuthPassword, ",")
		if len(users) != len(passwords) {
			panic("basic auth setting error")
		}
		for i := 0; i < len(users); i++ {
			accounts[users[i]] = passwords[i]
		}
		chat.POST("/chat-process", gin.BasicAuth(accounts), middlewares.RateLimitMiddleware(1, 2), chatService.ChatProcess)
	} else {
		chat.POST("/chat-process", middlewares.RateLimitMiddleware(1, 2), chatService.ChatProcess)
	}

	// chat.POST("/process", BasicAuth(accountService, r.OpsLink), middlewares.RateLimitMiddleware(1, 2), chatService.MessageProcess)
	chat.POST("/config", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "Success",
			"data": map[string]string{
				"apiModel":   "ChatGPTAPI",
				"socksProxy": "",
			},
		})
	})
	chat.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{})
	})
	chat.POST("/session", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "Success",
			"message": "",
			"data": gin.H{
				"auth": false,
			},
		})
	})
	// entry.POST("/accounts", OpsAuth(r.OpsKey), accountService.AccountProcess)
	// entry.Any("/admin/*relativePath", gin.BasicAuth(gin.Accounts{"admin": r.OpsKey}), func(ctx *gin.Context) {
	// 	if ctx.Request.URL.Path == "/admin/accounts" {
	// 		accountService.AccountProcess(ctx)
	// 	} else {
	// 		http.FileServer(http.Dir(path.Join(r.FrontendPath))).ServeHTTP(ctx.Writer, ctx.Request)
	// 	}
	// })

	// proxy.NoRoute(func(ctx *gin.Context) {
	// 	http.FileServer(http.Dir(r.config.Frontend.Path)).ServeHTTP(ctx.Writer, ctx.Request)
	// })

	// entry.Any("/webb/*path", func(ctx *gin.Context) {
	// 	ctx.Request.URL.Path = ctx.Param("path")
	// 	// proxy.ServeHTTP(ctx.Writer, ctx.Request)
	// 	http.FileServer(http.Dir(r.config.Frontend.Path)).ServeHTTP(ctx.Writer, ctx.Request)
	// })

	entry.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{})
	})

	server.Handler = entry
	go func(ctx context.Context) {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.InfoF("Server shutdown with error %v", err)
		}
	}(ctx)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.FatalF("Server listen and serve error %v", err)
	}
}

// func (r *ChatGPTWebServer) updateAssetsFiles() error {
// 	pairs := map[string]string{}
// 	// old := `{avatar:"https://raw.githubusercontent.com/Chanzhaoyu/chatgpt-web/main/src/assets/avatar.jpg",name:"ChenZhaoYu",description:'Star on <a href="https://github.com/Chanzhaoyu/chatgpt-bot" class="text-blue-500" target="_blank" >Github</a>'}`
// 	// new := fmt.Sprintf(`{avatar:"https://raw.githubusercontent.com/Chanzhaoyu/chatgpt-web/main/src/assets/avatar.jpg",name:"获取帮助输入/help",description:'<a href="%s" class="text-blue-500" target="_blank" >自助中心</a>'}`, link)
// 	// pairs[old] = new
// 	old := `{}.VITE_GLOB_OPEN_LONG_REPLY`
// 	new := `{VITE_GLOB_OPEN_LONG_REPLY:"true"}.VITE_GLOB_OPEN_LONG_REPLY`
// 	pairs[old] = new
// 	old = `<link rel="manifest" href="/manifest.webmanifest"><script id="vite-plugin-pwa:register-sw" src="/registerSW.js"></script>`
// 	new = ``
// 	pairs[old] = new
// 	old = `[y(" 此项目开源于 "),e("a",{class:"text-blue-600 dark:text-blue-500",href:"https://github.com/Chanzhaoyu/chatgpt-web",target:"_blank"}," Github "),y(" ，免费且基于 MIT 协议，没有任何形式的付费行为！ ")]`
// 	new = `[y(" 此项目开源于 "),e("a",{class:"text-blue-600 dark:text-blue-500",href:"https://chatgpt-web",target:"_blank"}," Github ")]`
// 	pairs[old] = new
// 	return utils.ReplaceFiles(r.config.Frontend.Path, pairs)
// }

func main() {
	loadDependOn()
	cnf := config.GetConf()

	// fmt.Printf("%+v\n", cnf)

	logger := log.NewLogger()
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer stop()
	app := NewChatGPTWebServer(config.GetConf(), logger)
	app.Run(ctx)
	<-ctx.Done()
}

func loadDependOn() {
	config.InitConf(cmd.Args.Config)
	cnf := config.GetConf()
	// fmt.Printf("%+v", cnf)
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetLevel(cnf.Log.Level)
	log.SetPrintCaller(true)
}
