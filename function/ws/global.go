package ws

import (
	"github.com/gorilla/websocket"
)

var Ip2Conn = map[string]*websocket.Conn{}
var Conn2ip = map[*websocket.Conn]string{}
var Conn2info = map[*websocket.Conn]Infomation{}

var MessageChan = make(chan Message, 20) //use for getting message
var PushChan = make(chan Message, 20)    //use for sending message

type Message struct {
	SubscribeType string
	Data          string
}

type Infomation struct {
	SubscribeTypes map[string]bool //订阅队列列表
}

//type Queue struct {
//	Clients       map[*client]bool
//	Broadcast     chan []byte
//	Quit          chan []byte
//	SubscribeType string //所属订阅队列
//}

//type Hub struct {
//	clients     map[*client]bool //这里类似我的user2conn
//	subscribe   chan *client
//	unsubscribe chan *client
//	register    chan *client
//	unregister  chan *client
//	UserQueue   []*Queue
//}

//type client struct {
//	hub            *Hub
//	conn           *websocket.Conn
//	send           chan []byte
//	Uid            string
//	SubscribeType  string   //订阅队列
//	SubscribeTypes []string //订阅队列列表
//}

type MsgPack struct {
	SocketType string `json:"socket_type"` //消息类型
	Subscribed int    `json:"subscribed"`  //1订阅 0取消订阅
}

//func NewHub() *Hub {
//	return &Hub{
//		subscribe:   make(chan *client),
//		unsubscribe: make(chan *client),
//		register:    make(chan *client),
//		unregister:  make(chan *client),
//		clients:     make(map[*client]bool),
//	}
//}

//func NewClient() *client {
//	return &client{
//		send:           make(chan []byte, maxMessageSize),
//		Uid:            Uuid(),
//		SubscribeType:  "",
//		SubscribeTypes: []string{},
//	}
//}

//func NewQueue() *Queue {
//	return &Queue{
//		Clients:   make(map[*client]bool),
//		Broadcast: make(chan []byte),
//		Quit:      make(chan []byte),
//	}
//}
