package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

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
	info.SubscribeTypes = make(map[string]bool)
	Conn2info[conn] = info

}

func On_close(conn *websocket.Conn) {
	On_exit(conn)
	// 发送 websocket 结束包
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	// 真正关闭 conn
	conn.Close()
}

func On_exit(conn *websocket.Conn) {
	delete(Ip2Conn, Conn2ip[conn])
	delete(Conn2info, conn)
	delete(Conn2ip, conn)
}
