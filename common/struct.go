/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package common

import (
	"github.com/gorilla/websocket"
	"sync"
)

////////// to client /////////////
// login
type LoginReq struct {
	Id   string `json:"id"`
	Cmd  string `json:"cmd"` // login
	Uuid uint64 `json:"uuid"`
}

// heartbeat
type HeartBeatReq struct {
	Id  string `json:"id"`
	Cmd string `json:"cmd"` // heartbeat
}

// push
type PushClientReq struct {
	Id  string `json:"id"`
	Cmd string `json:"cmd"`
	Msg string `json:"msg"`
}

type Response struct {
	Id   string `json:"id"`
	Cmd  string `json:"cmd"`
	Data struct {
		Code   uint32 `json:"code"`
		Result string `json:"result"`
	}
}

////////// to server /////////////
// push msg
type PushReq struct {
	Id   string
	Uuid uint64
	Msg  string
}

type PushRsp struct {
	Seq string
}

// query isOnline
type QueryReq struct {
	Uuid uint64
}

type QueryRsp struct {
	IsOnline bool
}

// route
type RouteReq struct {
	Id   string
	Uuid uint64
	Msg  string
}

type RouteRsp struct {
	Seq string
}

// push transfer
type TransferData struct {
	Seq  string `json:"seq"`
	Uuid uint64 `json:"uuid"`
	Msg  string `json:"msg"`
}

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
