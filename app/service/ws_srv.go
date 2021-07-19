/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/chuck1024/gd"
	"github.com/gorilla/websocket"
	"hydra/app/domain"
	"hydra/app/service/sp"
	"sync"
)

type Client struct {
	Id     string
	Socket *websocket.Conn
	Uuid   uint64
}

type ClientHub struct {
	Clients    map[uint64]*Client
	hLock      sync.RWMutex
	SendMsg    chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func (h *ClientHub) SetConnValue(uuid uint64, client *Client) {
	h.hLock.Lock()
	defer h.hLock.Unlock()

	h.Clients[uuid] = client
}

func (h *ClientHub) GetConnValue(uuid uint64) *Client {
	h.hLock.Lock()
	defer h.hLock.Unlock()

	if conn, ok := h.Clients[uuid]; ok {
		return conn
	}

	return nil
}

func (h *ClientHub) DelConnValue(uuid uint64) {
	h.hLock.Lock()
	defer h.hLock.Unlock()

	if _, ok := h.Clients[uuid]; ok {
		delete(h.Clients, uuid)
	}

	return
}

var Hub = ClientHub{
	SendMsg:    make(chan []byte, 10000),
	Register:   make(chan *Client, 10000),
	Unregister: make(chan *Client, 10000),
	Clients:    make(map[uint64]*Client),
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
					sp.Get().UidCache.DelUuid(client.Uuid)
				}
			}

		case sendMsg := <-Hub.SendMsg:
			data := &domain.TransferData{}
			err := json.Unmarshal(sendMsg, data)
			if err != nil {
				gd.Error("[Start] json unmarshal occur error:%s", err)
				continue
			}

			conn := Hub.GetConnValue(data.Uuid)
			if conn == nil {
				continue
			}

			pd := &domain.PushClientReq{
				Id:  data.Seq,
				Cmd: domain.PushCmd,
				Msg: data.Msg,
			}

			pdb, _ := json.Marshal(pd)

			conn.Socket.WriteMessage(websocket.TextMessage, pdb)
		}
	}
}
