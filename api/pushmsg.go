package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"main.go/function/ws"
)

func PushmsgArray(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("content-type", "application/json")
	body, _ := c.GetRawData()
	message := string(body)
	//fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), r.URL, message)
	data := gjson.Get(message, "data")
	if !data.Exists() || !data.IsArray() {
		fmt.Println("推送数据包data非法")
		c.String(200, "error")
		return
	}
	go func(data gjson.Result) {
		for _, result := range data.Array() {
			if !result.Get("socket_type").Exists() {
				continue
			}
			var msg ws.Message
			msg.SubscribeType = result.Get("socket_type").String()
			msg.Data = result.String()
			ws.MessageChan <- msg
		}
	}(data)
	c.String(200, "ok")
	return
}
