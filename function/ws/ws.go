package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"main.go/config"
	"main.go/tuuz/Date"
	"sync"
	"time"
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

func Ws_connect(hub *Hub, c *gin.Context) {
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
		ws_handler(hub, conn)
	}
}

func ws_handler(hub *Hub, conn *websocket.Conn) {
	//defer On_close(conn)
	////连入时发送欢迎消息
	//go On_connect(conn)
	//for {
	//	mt, d, err := conn.ReadMessage()
	//	conn.RemoteAddr()
	//	if mt == -1 {
	//		break
	//	}
	//	if err != nil {
	//		fmt.Println(mt)
	//		fmt.Printf("read fail = %v\n", err)
	//		break
	//	}
	//	Handler(string(d), conn)
	//}

	//ident here
	fmt.Println(conn.RemoteAddr())

	c := NewClient()
	c.hub = hub
	c.conn = conn

	hub.register <- c
	go c.writePump()
	c.readPump()
}

func On_connect(conn *websocket.Conn) {
	//err := conn.WriteMessage(1, []byte("连入成功"))
	message := map[string]interface{}{
		"remote_addr":  conn.RemoteAddr(),
		"connect_time": Date.Int2Date(time.Now().Unix()),
	}
	str := map[string]interface{}{
		"code": 0,
		"data": message,
		"type": "connected",
	}
	err := conn.WriteJSON(str)

	if err != nil {
		fmt.Printf("write fail = %v\n", err)
		return
	}
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

func Handler(json_str string, conn *websocket.Conn) {
	if config.DEBUG {
		fmt.Println("json_ws:", json_str)
	}
	//json, jerr := Jsong.JObject(json_str)
	//if jerr != nil {
	//	fmt.Println("jsonerr", jerr)
	//	return
	//
	//}
	//if config.DEBUG_WS_REQ {
	//	fmt.Println("DEBUG_WS_REQ", json_str)
	//}
	//data, derr := Jsong.ParseObject(json["data"])
	//if derr != nil {
	//	fmt.Println("ws_derr:", derr)
	//	data = map[string]interface{}{}
	//}
	//Type := Calc.Any2String(json["type"])
	//switch Type {
	//
	//default:
	//	fmt.Println("undefine_type:", Type)
	//	break
	//}
}
