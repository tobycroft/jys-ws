package cron

import "main.go/function/ws"

func Message_recv() {
	for message := range ws.MessageChan {
		for conn, infomation := range ws.Conn2info {
			if infomation.SubscribeTypes[message.SubscribeType] != nil {

			}
		}
	}
}

func Message_send() {

}
