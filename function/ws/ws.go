package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	http2 "net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http2.Request) bool {
		return true
	},
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
	//conn.WriteMessage(1, []byte("连入成功"))
	//ident here
	remoteaddr := conn.RemoteAddr().String()
	fmt.Println("远程连入：", remoteaddr)

	//Ip2Conn.Store(remoteaddr, conn)
	Conn2ip.Store(conn, remoteaddr)

	var info Infomation
	info.SubscribeTypes = make(map[string]bool)
	Conn2info.Store(conn, info)
	go UserMsgChan(conn)
}

func On_close(conn *websocket.Conn) {
	On_exit(conn)
	// 发送 websocket 结束包
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	// 真正关闭 conn
	conn.Close()
}

func On_exit(conn *websocket.Conn) {
	Conn2info.Delete(conn)
	Conn2ip.Delete(conn)
	ccc, has := Conn2Chan.Load(conn)
	if has {
		ccc.(chan string) <- "123"
	}
	Conn2Chan.Delete(conn)
	//ip, has := Conn2ip.LoadAndDelete(conn)
	//if has {
	//	Ip2Conn.Delete(ip)
	//}
}
