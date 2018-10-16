/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	"github.com/chuck1024/godog"
	"github.com/gorilla/websocket"
	"hydra/model"
	"hydra/service"
	"net/http"
)

func WsControl(resp http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(resp, req, nil)
	if err != nil {
		http.NotFound(resp, req)
		return
	}

	godog.Debug("[WsControl] %s connected.", conn.RemoteAddr().String())

	client := &model.Client{Id: conn.RemoteAddr().String(), Socket: conn}

	for {
		_, message, err := client.Socket.ReadMessage()
		if err != nil {
			service.Hub.Unregister <- client
			client.Socket.Close()
			godog.Debug("[WsControl] %s disconnected.", client.Socket.RemoteAddr().String())
			break
		}

		service.HandleData(message, client)
	}
}