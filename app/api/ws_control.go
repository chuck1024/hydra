/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package api

import (
	"github.com/chuck1024/gd"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"hydra/app/service"
	"net/http"
)

func WsControl(c *gin.Context) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	gd.Debug("[WsControl] %s connected.", conn.RemoteAddr().String())

	client := &service.Client{Id: conn.RemoteAddr().String(), Socket: conn}

	for {
		_, message, err := client.Socket.ReadMessage()
		if err != nil {
			service.Hub.Unregister <- client
			client.Socket.Close()
			gd.Debug("[WsControl] %s disconnected.", client.Socket.RemoteAddr().String())
			break
		}

		service.HandleData(message, client)
	}
}
