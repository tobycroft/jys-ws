package cron

import (
	"github.com/gorilla/websocket"
	"main.go/function/ws"
)

func Message_recv() {
	for message := range ws.MessageChan {
		ws.Conn2info.Range(func(conn, infomation interface{}) bool {
			if infomation.(ws.Infomation).SubscribeTypes[message.SubscribeType] == true {
				var psh ws.Push
				psh.Conn = conn.(*websocket.Conn)
				psh.Data = message.Data
				ws.PushChan <- psh
			}
			return true
		})
	}
}

func Message_send() {
	for push := range ws.PushChan {
		info, has := ws.Conn2info.Load(push.Conn)
		if has {
			go func(data string) {
				info.(ws.Infomation).Lock.Lock()
				push.Conn.WriteMessage(1, []byte(data))
				info.(ws.Infomation).Lock.Unlock()
			}(push.Data)

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
