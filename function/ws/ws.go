package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"main.go/config"
	http2 "net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http2.Request) bool {
		return true
	},
}

var upgrader_compress = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http2.Request) bool {
		return true
	},
	EnableCompression: true,
}

func Ws_connect(c *gin.Context, compress bool) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("content-type", "application/json")
	if !websocket.IsWebSocketUpgrade(c.Request) {
		return
	} else {
		if compress {
			conn, err := upgrader_compress.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				fmt.Printf("err = %s\n", err)
				return
			}
			ws_handler(conn, compress)
		} else {
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				fmt.Printf("err = %s\n", err)
				return
			}
			ws_handler(conn, compress)
		}

	}
}

func ws_handler(conn *websocket.Conn, compress bool) {
	defer On_close(conn)
	//连入时发送欢迎消息
	On_connect(conn)
	for {
		mt, d, err := conn.ReadMessage()
		//conn.RemoteAddr()
		if mt == -1 {
			return
		}
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			log.Printf("error: %v", err)
			return
		}
		if err != nil {
			fmt.Println(mt)
			fmt.Printf("read fail = %v\n", err)
			return
		}
		Handler(string(d), conn)
	}

}

func On_connect(conn *websocket.Conn) {
	//conn.WriteMessage(1, []byte("连入成功"))
	//ident here
	remoteaddr := conn.RemoteAddr().String()
	if config.DEBUG {
		fmt.Println("远程连入：", remoteaddr)
	}
	//Ip2Conn.Store(remoteaddr, conn)
	Conn2ip.Store(conn, remoteaddr)

	var info Infomation
	info.SubscribeTypes = new(sync.Map)

	Conn2info.Store(conn, info)
	Conn2Chan.Store(conn, make(chan string, 100))
	go UserMsgChan(conn)
}

func On_close(conn *websocket.Conn) {
	// 发送 websocket 结束包
	//conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	// 真正关闭 conn

	On_exit(conn)
	conn.Close()
}

func On_exit(conn *websocket.Conn) {
	ccc, has := Conn2Chan.Load(conn)
	if has {
		timeout := time.NewTimer(time.Microsecond * 500)
		select {
		case ccc.(chan string) <- "close":
			break
		case <-timeout.C:
			break
		}
		Conn2Chan.Delete(conn)
	}
	Conn2info.Delete(conn)
	Conn2ip.Delete(conn)
}
