package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
)

func UserMsgChan(conn *websocket.Conn) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()
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
