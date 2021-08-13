package ws

import "github.com/gorilla/websocket"

func UserMsgChan(conn *websocket.Conn) {
	for data := range Conn2Chan {
		_, has := Conn2info.Load(conn)
		if !has {
			return
		}

		conn.WriteMessage(1, []byte(data))
	}
}
