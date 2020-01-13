/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	"github.com/chuck1024/doglog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"hydra/common"
	"hydra/service"
	"net/http"
)

func WsControl(c *gin.Context) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	doglog.Debug("[WsControl] %s connected.", conn.RemoteAddr().String())

	client := &common.Client{Id: conn.RemoteAddr().String(), Socket: conn}

	for {
		_, message, err := client.Socket.ReadMessage()
		if err != nil {
			service.Hub.Unregister <- client
			client.Socket.Close()
			doglog.Debug("[WsControl] %s disconnected.", client.Socket.RemoteAddr().String())
			break
		}

		service.HandleData(message, client)
	}
}
