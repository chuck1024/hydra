/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package model

import (
	"github.com/gorilla/websocket"
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

func (h *ClientHub) DelConnValue(uuid uint64) error {
	h.hLock.Lock()
	defer h.hLock.Unlock()

	if _, ok := h.Clients[uuid]; ok {
		delete(h.Clients, uuid)
	}

	return nil
}
