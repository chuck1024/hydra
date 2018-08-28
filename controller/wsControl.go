/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	"github.com/chuck1024/godog"
	"github.com/gorilla/websocket"
	"net/http"
	"hydra/common"
	"hydra/service/core"
	"hydra/cache"
)

func WsControl(resp http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(resp, req, nil)
	if err != nil {
		http.NotFound(resp, req)
		return
	}

	godog.Debug("[WsControl] %s connected.", conn.RemoteAddr().String())

	client := &common.Client{Id: conn.RemoteAddr().String(), Socket: conn}

	for {
		_, message, err := client.Socket.ReadMessage()
		if err != nil {
			core.Hub.Unregister <- client
			client.Socket.Close()
			godog.Debug("[WsControl] %s disconnected.", client.Socket.RemoteAddr().String())
			break
		}

		dealConn(message, client)
	}
}

func dealConn(message []byte, client *common.Client){
	//handle message according to yourself
	godog.Debug("[dealConn] receive message:%s",string(message))

	err := cache.SetUuid(client.Uuid)
	if err != nil {
		godog.Error("[dealConn] cache set uuid occur error:%s",err)
		return
	}
}