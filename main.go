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
	hub := ws.NewHub()
	go hub.Run()

	timelocal, _ := time.LoadLocation("Asia/Chongqing")
	time.Local = timelocal

	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	gin.DefaultWriter = ioutil.Discard
	// websocket echo

	r.GET("/", func(c *gin.Context) {
		r := c.Request
		w := c.Writer
		ws.Ws_connect(hub, c, w, r)
	})

	r.GET("/ws", func(c *gin.Context) {
		ws.Ws_connect(hub, c)
	})

	r.POST("/pushmsg", func(c *gin.Context) {
		api.Pushmsg(hub, c) //推送消息
	})

	r.POST("/pushmsgarray", func(c *gin.Context) {
		api.PushmsgArray(hub, c) //推送消息
	})

	r.GET("/wsjava", func(c *gin.Context) {
		r := c.Request
		w := c.Writer
		ws.ServeJavaWs(hub, w, r)
	})

	go r.Run(config.SERVER_LISTEN_ADDR + ":" + config.SERVER_LISTEN_PORT1)
	go r.Run(config.SERVER_LISTEN_ADDR + ":" + config.SERVER_LISTEN_PORT2)
	go r.RunTLS(config.SERVER_LISTEN_ADDR+":"+config.SERVER_LISTEN_PORT_SSL, "cert.pem", "cert.key")

	if err := http2.ListenAndServe("0.0.0.0:"+config.SERVER_DEBUG_PORT, nil); err != nil {
		fmt.Printf("start pprof failed on %s\n", config.SERVER_DEBUG_PORT)
	}
}

func ws_handler(conn *websocket.Conn) {
	defer ws.On_close(conn)
	//连入时发送欢迎消息
	go ws.On_connect(conn)
	for {
		mt, d, err := conn.ReadMessage()
		conn.RemoteAddr()
		if mt == -1 {
			break
		}
		if err != nil {
			fmt.Println(mt)
			fmt.Printf("read fail = %v\n", err)
			break
		}
		ws.Handler(string(d), conn)
	}
}
