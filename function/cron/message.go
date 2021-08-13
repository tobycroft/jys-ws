package cron

import (
	"main.go/function/ws"
)

func Message_recv() {
	for message := range ws.MessageChan {
		for conn, infomation := range ws.Conn2info {
			//		if infomation.SubscribeTypes[message.SubscribeType] == true {
			var psh ws.Push
			psh.Conn = conn
			psh.Data = message.Data
			ws.PushChan <- psh
			//}
		}
	}
}

func Message_send() {
	for push := range ws.PushChan {
		if ws.Conn2ip[push.Conn] != "" {
			push.Conn.WriteJSON(push.Data)
			//push.Conn.(*websocket.Conn).WriteJSON(push.Data)
		}
	}
}

//func socket_send_handle(uid string, channel chan interface{}) {
//	for message := range channel {
//		conn, has := User2Conn2.Load(uid)
//		if has {
//			conn.(*websocket.Conn).WriteJSON(message)
//		} else {
//			return
//		}
//	}
//}
