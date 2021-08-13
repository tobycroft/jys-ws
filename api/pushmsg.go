package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"main.go/function/ws"
)

type msg_struct struct {
	SocketType string `json:"socket_type"`
}

func Pushmsg(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("content-type", "application/json")
	body, _ := c.GetRawData()
	//fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), body)
	var data msg_struct
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("拆解推送数据包失败:", err.Error())
		c.String(200, "error")
		return
	}
	go func(data msg_struct, json string) {
		if data.SocketType == "" {
			fmt.Println("推送数据包没有类型")
			c.String(200, "error")
			return
		}
		var msg ws.Message
		msg.SubscribeType = data.SocketType
		msg.Data = json

		ws.MessageChan <- msg
	}(data, string(body))
	c.String(200, "ok")
	return
}

func PushmsgArray(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("content-type", "application/json")
	body, _ := c.GetRawData()
	message := string(body)

	data := gjson.Get(message, "data")
	if !data.Exists() || !data.IsArray() {
		fmt.Println("推送数据包data非法", message)
		c.String(200, "error")
		return
	}
	//fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), message)
	go func(data gjson.Result) {
		//defer func() {
		//	if r := recover(); r != nil {
		//		fmt.Printf("捕获到的错误：%s\n", r)
		//	}
		//}()
		//根据消息类型，向指定订阅队列发广播
		ws.Conn2info.Range(func(conn, infomation interface{}) bool {
			for _, result := range data.Array() {
				if !result.Get("socket_type").Exists() {
					continue
				}
				if infomation.(ws.Infomation).SubscribeTypes[result.Get("socket_type").String()] == true {
					var push ws.Push
					push.Conn = conn.(*websocket.Conn)
					push.Data = result.String()
					ws.PushChan <- push
				}

			}
			return true
		})
	}(data)
	c.String(200, "ok")
	return
}
