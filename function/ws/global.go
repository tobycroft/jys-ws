package ws

import (
	"github.com/gorilla/websocket"
	"sync"
)

//var Ip2Conn = map[string]*websocket.Conn{}
//var Ip2Conn sync.Map

//var Conn2ip = map[*websocket.Conn]string{}
var Conn2ip sync.Map

//var Conn2Chan = map[*websocket.Conn]chan string{}
var Conn2Chan sync.Map

//var Conn2info = map[*websocket.Conn]Infomation{}
var Conn2info sync.Map

var MessageChan = make(chan Message, 100) //use for getting message
var PushChan = make(chan Push, 100)       //use for sending message

type Message struct {
	SubscribeType string
	Data          string
}

type Push struct {
	Conn *websocket.Conn
	Data string
}

type Infomation struct {
	//SubscribeTypes map[string]bool //订阅队列列表
	SubscribeTypes *sync.Map //订阅队列列表
}

type MsgPack struct {
	SocketType string `json:"socket_type"` //消息类型
	Subscribed int    `json:"subscribed"`  //1订阅 0取消订阅
}
