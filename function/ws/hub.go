package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,

	//CheckOrigin:     checkOrigin,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), "建立连接:", client.Uid)
		case client := <-h.subscribe:
			fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), "订阅队列:", client.Uid, client.SubscribeType)
			//判断队列是否为空，或者队列平均人数大于100人 增加队列   || float32(len(h.clients)) / float32(len(h.UserQueue)) > 100
			if len(h.UserQueue) == 0 {
				queue := NewQueue()
				queue.Clients[client] = true
				queue.SubscribeType = client.SubscribeType
				h.UserQueue = append(h.UserQueue, queue)
				go queue.Ws_handler()
				fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), "队列创建:", len(h.UserQueue))
			} else {
				queue := NewQueue()
				bFind := false
				for _, q := range h.UserQueue {
					//查找订阅队列是否存在 并筛选人数最少的队列
					fmt.Println("遍历队列", q.SubscribeType)
					if q.SubscribeType == client.SubscribeType {
						fmt.Println("找到队列,队列客户端数量:", len(q.Clients))
						bFind = true
						if len(queue.Clients) == 0 || len(q.Clients) < len(queue.Clients) {
							queue = q
						}
					}
				}
				//如果没有找到队列 或队列人数超过100，增加新队列
				if bFind == false || len(queue.Clients) >= 100 {
					if bFind == false {
						fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), "没有找到队列,新队列创建:", len(h.UserQueue))
					} else {

					}
					queue = NewQueue()
					queue.SubscribeType = client.SubscribeType
					h.UserQueue = append(h.UserQueue, queue)
					go queue.Ws_handler()
				}
				queue.Clients[client] = true
			}
			fmt.Printf("%s 总连接:%d,队列数:%d\n", time.Now().Local().Format("2006-01-02 15:04:05"), len(h.clients), len(h.UserQueue))
			/*
				fmt.Println("-----------------------------")
				for _,queue := range h.UserQueue{
					for c,_ := range queue.Clients{
						fmt.Printf("%s ",c.Uid)
					}
					fmt.Println("\n****************************")
				}
				fmt.Println("-----------------------------")
			*/
		case client := <-h.unsubscribe:
			fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), "取消订阅:", client.Uid, client.SubscribeType)
			for i, q := range h.UserQueue {
				//查找订阅队列把客户端移除
				fmt.Println("遍历队列", q.SubscribeType)
				if q.SubscribeType == client.SubscribeType {
					fmt.Println("找到队列,队列客户端数量:", len(q.Clients))
					_, ok := q.Clients[client]
					if ok {
						delete(q.Clients, client)
					}
					//队列客户端为空释放
					if len(q.Clients) == 0 {
						h.UserQueue = append(h.UserQueue[:i], h.UserQueue[i+1:]...)
						q.Quit <- []byte("quit")
						fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), "队列释放:", client.SubscribeType)
						q = nil
					}
				}
			}
			fmt.Printf("%s 总连接:%d,队列数:%d\n", time.Now().Local().Format("2006-01-02 15:04:05"), len(h.clients), len(h.UserQueue))
		case client := <-h.unregister:
			_, ok := h.clients[client]
			if ok {
				delete(h.clients, client)
				close(client.send)
				fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), "断开连接:", client.Uid)
				//遍历队列，从队列把客户端移除
				for i, queue := range h.UserQueue {
					_, ok := queue.Clients[client]
					if ok {
						delete(queue.Clients, client)
						if len(queue.Clients) == 0 {
							queue.Quit <- []byte("quit")
							fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), "队列释放:", queue.SubscribeType)
							queue = nil
							if len(h.UserQueue) > 1 {
								if len(h.UserQueue) >= i-1 {
									h.UserQueue = append(h.UserQueue[:i], h.UserQueue[i+1:]...)
								}
							} else {
								h.UserQueue = make([]*Queue, 0)
							}
						}
						//fmt.Println("queue.clients",queue.Clients)
						//break
					}
				}
				//fmt.Println("UserQueue",h.UserQueue)
				fmt.Printf("%s 总连接:%d,队列数:%d\n", time.Now().Local().Format("2006-01-02 15:04:05"), len(h.clients), len(h.UserQueue))
				/*
					fmt.Println("-----------------------------")
					for _,queue := range h.UserQueue{
						for c,_ := range queue.Clients{
							fmt.Printf("%s ",c.Uid)
						}
						fmt.Println("\n****************************")
					}
					fmt.Println("-----------------------------")
				*/
			}
			//case m := <-h.broadcast:
			//	Distribute(h,m)
		}
	}
}

func (s *Queue) Ws_handler() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
		close(s.Broadcast)
		close(s.Quit)

	}()
	for {
		select {
		case message := <-s.Broadcast:
			for c, _ := range s.Clients {
				c.send <- message
				//if err := c.write(websocket.TextMessage, message); err != nil {
				//	c.hub.unregister <- c
				//	c.conn.Close()
				//}
			}
		case <-s.Quit:
			return
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
		ticker.Stop()
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		select {
		case message := <-c.send:
			if err := c.write(websocket.TextMessage, message); err != nil {
				fmt.Println("发数据包失败:", err.Error())
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				fmt.Println("发心跳包失败:", err.Error())
				return
			}
		}
	}
}

func (c *client) write(mt int, message []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, message)
}

func ServeJavaWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	ac := NewClient()
	ac.hub = hub
	ac.conn = conn

	go ac.writePumpJava()
	ac.readPumpJava()
}

func (c *client) writePumpJava() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message := <-c.send:
			if err := c.write(websocket.TextMessage, message); err != nil {
				fmt.Println("发数据包失败:", err.Error())
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				fmt.Println("发心跳包失败:", err.Error())
				return
			}
		}
	}
}

func (c *client) readPumpJava() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, body, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		body = bytes.TrimSpace(bytes.Replace(body, newline, space, -1))
		if string(body) == "heart" {
			fmt.Println("收到心跳包")
			continue
		}
		fmt.Println("收到java_websocket消息:", string(body))

		p := NewMsgPack()
		err = json.Unmarshal(body, &p)
		if err != nil {
			fmt.Println("拆解推送数据包失败:", err.Error())
			c.send <- []byte("error")
			continue
		}
		if p.SocketType == "" {
			fmt.Println("推送数据包没有类型")
			c.send <- []byte("error")
			continue
		}

		//根据消息类型，向指定订阅队列发广播
		for _, queue := range c.hub.UserQueue {
			if queue.SubscribeType == p.SocketType {
				queue.Broadcast <- body
			}
		}

		c.send <- []byte("ok")
	}
}
