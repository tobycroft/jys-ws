package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"main.go/config"
	"main.go/tuuz/Calc"
	"main.go/tuuz/Date"
	"main.go/tuuz/Jsong"
	"net/http"
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

func Ws_connect(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.Method)
	if r.Method != "GET" {
		http.Error(w, "Method not allowd", 405)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	conn, err := upgrader.Upgrade(w, r, nil)
	//ident here
	fmt.Println(conn.RemoteAddr())
	if err != nil {
		log.Println(err)
		return
	}

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
	json, jerr := Jsong.JObject(json_str)
	if jerr != nil {
		fmt.Println("jsonerr", jerr)
		return

	}
	if config.DEBUG_WS_REQ {
		fmt.Println("DEBUG_WS_REQ", json_str)
	}
	data, derr := Jsong.ParseObject(json["data"])
	if derr != nil {
		fmt.Println("ws_derr:", derr)
		data = map[string]interface{}{}
	}
	Type := Calc.Any2String(json["type"])
	switch Type {
	case "init", "INIT":
		auth_init(conn, data, Type)
		break

	case "join_room", "JOIN_ROOM":
		join_room(conn, data, Type)
		break

	case "exit_room", "EXIT_ROOM":
		exit_room(conn, data, Type)
		break

	case "msg_list", "MSG_LIST":
		msg_list(conn, data, Type)
		break

	case "private_msg", "PRIVATE_MSG":
		private_msg(conn, data, Type)
		break

	case "group_msg":
		group_msg(conn, data, Type)
		break

	case "requst_count":
		requst_count(conn, data, Type)
		break

	case "ping":
		ping(conn, data, Type)
		break

	case "api":
		api(conn, data, Type)
		break

	case "clear_private_unread":
		clear_private_unread(conn, data, Type)
		break

	case "clear_group_unread":
		clear_group_unread(conn, data, Type)
		break

	default:
		fmt.Println("undefine_type:", Type)
		break
	}
}
