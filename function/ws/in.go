package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"main.go/config"
)

func Handler(json_str string, conn *websocket.Conn) {
	if config.DEBUG {
		fmt.Println("json_ws:", json_str)
	}
	if json_str == "heart" {
		//fmt.Println("收到心跳包")
		return
	}
	var msg MsgPack
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.UnmarshalFromString(json_str, &msg)
	if err != nil {
		fmt.Println("拆解客户端数据包失败:", err.Error())
		return
	}

	if msg.Subscribed == 0 {
		//info := Conn2info[conn]
		delete(Conn2info[conn].SubscribeTypes, msg.SocketType)
	} else {
		Conn2info[conn].SubscribeTypes[msg.SocketType] = true
	}
}
