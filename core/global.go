package core

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
)

type Queue struct {
	Clients map[*client]bool
	Broadcast	chan []byte
	Quit	chan []byte
	SubscribeType string   //所属订阅队列
}

type Hub struct {
	clients    map[*client]bool
	subscribe  chan *client
	unsubscribe chan *client
	register   chan *client
	unregister chan *client
	UserQueue  []*Queue
}

type client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	Uid string
	SubscribeType	string  //订阅队列
	SubscribeTypes []string  //订阅队列列表
}

type MsgPack struct {
	SocketType string `json:"socket_type"`  //消息类型
	Subscribed	int	`json:"subscribed"`  //1订阅 0取消订阅
}

func NewHub() *Hub {
	return &Hub{
		subscribe:  make(chan *client),
		unsubscribe:  make(chan *client),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
	}
}

func NewClient() *client{
	return &client{
		send:	make(chan []byte, maxMessageSize),
		Uid:	Uuid(),
		SubscribeType: "",
		SubscribeTypes: []string{},
	}
}

func NewQueue()  *Queue {
	return &Queue{
	Clients:    make(map[*client]bool),
	Broadcast: make(chan []byte),
	Quit:	make(chan []byte),
	}
}

func NewMsgPack()	*MsgPack {
	return &MsgPack{}
}

func Uuid()  string {
	u,_ := uuid.NewV4()
	return u.String()
}

func GetMd5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
