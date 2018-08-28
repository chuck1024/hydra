/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package core

import (
	"hydra/common"
	"github.com/gorilla/websocket"
	"hydra/cache"
	"encoding/json"
	"github.com/chuck1024/godog"
)

var Hub = common.ClientHub{
	SendMsg:    make(chan []byte, 10000),
	Register:   make(chan *common.Client, 10000),
	Unregister: make(chan *common.Client, 10000),
	Clients:    make(map[uint64]*common.Client),
}

func Start() {
	for {
		select {
		case client := <-Hub.Register:
			Hub.SetConnValue(client.Uuid, client)

		case client := <-Hub.Unregister:
			if client.Uuid > 0 {
				if _, ok := Hub.Clients[client.Uuid]; ok {
					Hub.DelConnValue(client.Uuid)
					cache.DelUuid(client.Uuid)
				}
			}

		case sendMsg := <-Hub.SendMsg:
			data := &common.TransferData{}
			err := json.Unmarshal(sendMsg,data)
			if err != nil {
				godog.Error("[Start] json unmarshal occur error:%s",err)
				continue
			}

			conn := Hub.GetConnValue(data.Uuid)
			if conn == nil {
				continue
			}

			conn.Socket.WriteMessage(websocket.TextMessage, []byte(data.Msg))
		}
	}
}
