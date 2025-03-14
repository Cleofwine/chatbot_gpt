package main

import (
	"chatgpt-qq/cmd/cqhttp"
	"chatgpt-qq/config"
	"chatgpt-qq/log"
	"context"
	"os"
	"os/signal"
)

func main() {
	cnf := config.GetConf()
	log.SetLevel(cnf.Log.Level)
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))

	// fmt.Printf("%+v\n", cnf)

	if cnf.CqHttp.WsServerPort != 0 {
		// 反向websocket，此应用作为cqhttp的服务端
		// 启动websocket server
		go cqhttp.RunWsServer()
		log.InfoF("监听 ws://%s:%d", cnf.CqHttp.WsServerHost, cnf.CqHttp.WsServerPort)
	} else {
		// 正向websocket，此应用作为cqhttp的客户端
		go cqhttp.Run()
	}

	// go cqhttp.Run()
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()
	<-ctx.Done()
}
