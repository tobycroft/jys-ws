package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
	"main.go/config"
)

func readPump() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if string(message) == "heart" {
			//fmt.Println("收到心跳包")
			continue
		}
		//fmt.Println("收到websocket消息:", string(message))
		//解析订阅 取消订阅
		p := NewMsgPack()
		err = json.Unmarshal(message, &p)
		if err != nil {
			fmt.Println("拆解客户端数据包失败:", err.Error())
			continue
		}
		c.SubscribeType = p.SocketType
		if p.Subscribed == 0 {
			c.hub.unsubscribe <- c
		} else {
			c.hub.subscribe <- c
		}
	}
}

func Handler(json_str string, conn *websocket.Conn) {
	if config.DEBUG {
		fmt.Println("json_ws:", json_str)
	}
	var msg MsgPack
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.UnmarshalFromString(json_str, &msg)
	if err != nil {
		fmt.Println("拆解客户端数据包失败:", err.Error())
		return
	}

	if msg.Subscribed == 0 {
		info := Conn2info[conn]
		info.SubscribeType = ""

	} else {
		info := Conn2info[conn]
		info.SubscribeType = msg.SocketType
		info.SubscribeTypes[msg.SocketType] = true
		Conn2info[conn] = info
	}
}
