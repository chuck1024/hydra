/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package common

import (
	"github.com/gorilla/websocket"
	"sync"
)

// 推送消息
type PushReq struct {
	Id   string
	Uuid uint64
	Msg  string
}

type PushRsp struct{
	Seq string
}

// 查询是否在线
type QueryReq struct {
	Uuid uint64
}

type QueryRsp struct{
	IsOnline bool
}

// route到链接所在实例
type RouteReq struct {
	Uuid uint64
	Msg string
}

type RouteRsp struct {
	Seq string
}

// push数据中间层
type TransferData struct {
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