package main

import (
	"fmt"
	"main.go/api"
	"main.go/config"
	"main.go/core"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {

	/* 创建集合 */
	hub := core.NewHub()
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
		core.ServeWs(w, r)
	})
	http.HandleFunc("/wsjava", func(w http.ResponseWriter, r *http.Request) {
		core.ServeJavaWs(hub, w, r)
	})
	fmt.Println("开始监听:", config.SERVER_LISTEN_PORT)
	go http.ListenAndServe(":"+config.SERVER_LISTEN_PORT, nil)

	if err := http.ListenAndServe(":"+config.SERVER_DEBUG_PORT, nil); err != nil {
		fmt.Printf("start pprof failed on %s\n", config.SERVER_DEBUG_PORT)
	}
}
