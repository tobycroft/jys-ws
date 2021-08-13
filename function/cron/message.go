package cron

import (
	"github.com/gorilla/websocket"
	"main.go/function/ws"
	"time"
)

func Message_recv() {
	for message := range ws.MessageChan {
		ws.Conn2info.Range(func(conn, infomation interface{}) bool {
			bo, has := infomation.(ws.Infomation).SubscribeTypes.Load(message.SubscribeType)
			if has && bo == true {
				var psh ws.Push
				psh.Conn = conn.(*websocket.Conn)
				psh.Data = message.Data

				timeout := time.NewTimer(time.Microsecond * 500)
				select {
				case ws.PushChan <- psh:
					break
				case <-timeout.C:
					break
				}
			}
			return true
		})
	}
}

func Message_send() {
	for push := range ws.PushChan {
		_, has := ws.Conn2info.Load(push.Conn)
		if has {
			cc, has := ws.Conn2Chan.Load(push.Conn)
			if has {
				timeout := time.NewTimer(time.Microsecond * 500)
				select {
				case cc.(chan string) <- push.Data:
					break
				case <-timeout.C:
					break
				}
			}
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
