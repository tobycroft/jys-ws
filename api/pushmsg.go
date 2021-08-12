package api

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"main.go/ws"
	"net/http"
)

func Pushmsg(h *ws.Hub, w http.ResponseWriter, r *http.Request) {
	//func  ( h *core.Hub, w *http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	//message := string(body)
	//fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"),r.URL,message)
	p := ws.NewMsgPack()
	err := json.Unmarshal(body, &p)
	if err != nil {
		fmt.Println("拆解推送数据包失败:", err.Error())
		w.Write([]byte("error"))
		return
	}
	if p.SocketType == "" {
		fmt.Println("推送数据包没有类型")
		w.Write([]byte("error"))
		return
	}
	//go func() {
	//	//根据消息类型，向指定订阅队列发广播
	//	for _, queue := range h.UserQueue {
	//		if queue.SubscribeType == p.SocketType {
	//			queue.Broadcast <- []byte(message)
	//		}
	//	}
	//}()
	w.Write([]byte("ok"))
	return
}

func PushmsgArray(h *ws.Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	message := string(body)
	//fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), r.URL, message)
	data := gjson.Get(message, "data")
	if !data.Exists() || !data.IsArray() {
		fmt.Println("推送数据包data非法")
		w.Write([]byte("error"))
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("捕获到的错误：%s\n", r)
			}
		}()
		//根据消息类型，向指定订阅队列发广播
		for _, queue := range h.UserQueue {
			for i := 0; i < len(data.Array()); i++ {
				if !data.Array()[i].Get("socket_type").Exists() {
					continue
				}
				if queue.SubscribeType == data.Array()[i].Get("socket_type").String() {
					msg := data.Array()[i].String()
					//fmt.Println(msg)
					queue.Broadcast <- []byte(msg)
				}
			}
		}
	}()

	w.Write([]byte("ok"))
	return
}
