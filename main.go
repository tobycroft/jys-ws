package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"main.go/api"
	"main.go/config"
	"main.go/ws"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {

	/* 创建集合 */
	hub := ws.NewHub()
	go hub.Run()
	timelocal, _ := time.LoadLocation("Asia/Chongqing")
	time.Local = timelocal

	http.HandleFunc("/pushmsg", func(w http.ResponseWriter, r *http.Request) {
		api.Pushmsg(hub, w, r) //推送消息
	})
	http.HandleFunc("/pushmsgarray", func(w http.ResponseWriter, r *http.Request) {
		api.PushmsgArray(hub, w, r) //推送消息
	})
	//http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.Ws_connect(w, r)
	})
	http.HandleFunc("/wsjava", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeJavaWs(hub, w, r)
	})
	fmt.Println("开始监听:", config.SERVER_LISTEN_PORT)
	go http.ListenAndServe(":"+config.SERVER_LISTEN_PORT, nil)

	if err := http.ListenAndServe("0.0.0.0:"+config.SERVER_DEBUG_PORT, nil); err != nil {
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
