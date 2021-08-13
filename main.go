package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"main.go/api"
	"main.go/config"
	"main.go/function/ws"
	http2 "net/http"
	_ "net/http/pprof"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http2.Request) bool {
		return true
	},
}

func main() {

	/* 创建集合 */
	go ws.Run()

	timelocal, _ := time.LoadLocation("Asia/Chongqing")
	time.Local = timelocal

	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	gin.DefaultWriter = ioutil.Discard
	// websocket echo

	//r.GET("/", func(c *gin.Context) {
	//	ws.Ws_connect(hub, c)
	//})

	r.GET("/ws", func(c *gin.Context) {
		ws.Ws_connect(c)
	})

	r.POST("/pushmsg", func(c *gin.Context) {
		api.Pushmsg(c) //推送消息
	})

	r.POST("/pushmsgarray", func(c *gin.Context) {
		api.PushmsgArray(c) //推送消息
	})

	r.GET("/wsjava", func(c *gin.Context) {
		ws.ServeJavaWs(c)
	})

	go r.Run(config.SERVER_LISTEN_ADDR + ":" + config.SERVER_LISTEN_PORT1)
	go r.Run(config.SERVER_LISTEN_ADDR + ":" + config.SERVER_LISTEN_PORT2)
	go r.RunTLS(config.SERVER_LISTEN_ADDR+":"+config.SERVER_LISTEN_PORT_SSL, "cert.pem", "cert.key")

	if err := http2.ListenAndServe("0.0.0.0:"+config.SERVER_DEBUG_PORT, nil); err != nil {
		fmt.Printf("start pprof failed on %s\n", config.SERVER_DEBUG_PORT)
	}
}
