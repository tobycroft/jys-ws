package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var Conn2User = make(map[*websocket.Conn]string)

var User2Conn2 sync.Map
var Conn2User2 sync.Map
var Room2 sync.Map

var User2Chan2 sync.Map

func socket_send_handle(uid string, channel chan interface{}) {
	for message := range channel {
		conn, has := User2Conn2.Load(uid)
		if has {
			conn.(*websocket.Conn).WriteJSON(message)
		} else {
			return
		}
	}
}

func Ws_connect(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("content-type", "application/json")
	if !websocket.IsWebSocketUpgrade(c.Request) {
		return
	} else {

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Printf("err = %s\n", err)
			return
		}
		ws_handler(conn)
	}
}

func ws_handler(conn *websocket.Conn) {
	defer On_close(conn)
	//连入时发送欢迎消息
	On_connect(conn)
	for {
		mt, d, err := conn.ReadMessage()
		conn.RemoteAddr()
		if mt == -1 {
			break
		}
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			log.Printf("error: %v", err)
			break
		}
		if err != nil {
			fmt.Println(mt)
			fmt.Printf("read fail = %v\n", err)
			break
		}
		Handler(string(d), conn)
	}

}

func On_connect(conn *websocket.Conn) {
	conn.WriteMessage(1, []byte("连入成功"))
	//ident here
	remoteaddr := conn.RemoteAddr().String()
	fmt.Println("远程连入：", remoteaddr)

	Ip2Conn[remoteaddr] = conn
	Conn2ip[conn] = remoteaddr

	var info Infomation
	Conn2info[conn] = info

	//go c.writePump()
	//c.readPump()
}

func On_close(conn *websocket.Conn) {
	On_exit(conn)
	// 发送 websocket 结束包
	conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 真正关闭 conn
	conn.Close()
}

func On_exit(conn *websocket.Conn) {
	uid, has := Conn2User2.Load(conn)
	//删除用户数据
	if has {
		Room2.Delete(uid.(string))
		User2Conn2.Delete(uid.(string))
		ch, has := User2Chan2.Load(uid.(string))
		if has {
			ch.(chan interface{}) <- "close"
		}
		User2Chan2.Delete(uid.(string))
		Conn2User2.Delete(conn)
	}
}
