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
		info, has := Conn2info.Load(conn)
		if has {
			info.(Infomation).SubscribeTypes.Delete(msg.SocketType)
			//delete(info.(Infomation).SubscribeTypes, msg.SocketType)
			Conn2info.Store(conn, info)
		}
	} else {
		info, has := Conn2info.Load(conn)
		if has {
			//替换模式，只能关注一种数据的模式
			//info.(Infomation).SubscribeTypes.Range(func(key, value interface{}) bool {
			//	info.(Infomation).SubscribeTypes.Delete(key)
			//	return true
			//})

			//同时关注模式
			info.(Infomation).SubscribeTypes.Store(msg.SocketType, true)
			//info.(Infomation).SubscribeTypes[msg.SocketType] = true
			Conn2info.Store(conn, info)
		}
	}
}
