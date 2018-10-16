/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/chuck1024/godog"
	"github.com/gorilla/websocket"
	"hydra/dao/cache"
	"hydra/model"
)

var Hub = model.ClientHub{
	SendMsg:    make(chan []byte, 10000),
	Register:   make(chan *model.Client, 10000),
	Unregister: make(chan *model.Client, 10000),
	Clients:    make(map[uint64]*model.Client),
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
			data := &model.TransferData{}
			err := json.Unmarshal(sendMsg, data)
			if err != nil {
				godog.Error("[Start] json unmarshal occur error:%s", err)
				continue
			}

			conn := Hub.GetConnValue(data.Uuid)
			if conn == nil {
				continue
			}

			pd := &model.PushClientReq{
				Id:  data.Seq,
				Cmd: "push",
				Msg: data.Msg,
			}

			pdb, _ := json.Marshal(pd)

			conn.Socket.WriteMessage(websocket.TextMessage, pdb)
		}
	}
}