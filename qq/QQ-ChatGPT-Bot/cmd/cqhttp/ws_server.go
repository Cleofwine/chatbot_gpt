package cqhttp

import (
	"chatgpt-qq/config"
	"chatgpt-qq/log"
	"fmt"
	"os"

	"net/http"

	"github.com/gorilla/websocket"
)

func RunWsServer() {
	logger := log.NewLogger()
	logger.SetOutput(os.Stderr)
	logger.SetLevel("info")
	http.Handle("/ws", new(wsHandler))
	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})
	cnf := config.GetConf()
	// 1. 首先建立一个http连接
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", cnf.CqHttp.WsServerHost, cnf.CqHttp.WsServerPort), nil)
	if err != nil {
		logger.Fatal(err)
	}
}

type wsHandler struct {
}

func (*wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accessToken := r.FormValue("ws_access_token")
	cnf := config.GetConf()
	if accessToken != cnf.CqHttp.WsAccessToken {
		log.Error("权限校验失败")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var upgrader = websocket.Upgrader{}
	// 2. 升级http连接为websocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()
	bot := NewBot()
	bot.Conn = conn
	bot.Read()
}
