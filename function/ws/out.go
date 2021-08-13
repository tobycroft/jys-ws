package ws

import "github.com/gorilla/websocket"

func UserMsgChan(conn *websocket.Conn) {
	cc, has := Conn2Chan.Load(conn)
	if has {
		for data := range cc.(chan string) {
			_, has := Conn2info.Load(conn)
			if !has {
				return
			}

			conn.WriteMessage(1, []byte(data))
		}
	}

}
