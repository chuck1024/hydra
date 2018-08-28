/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package core

import (
	"encoding/json"
	"github.com/chuck1024/godog"
	de "github.com/chuck1024/godog/error"
	"github.com/gorilla/websocket"
	"hydra/cache"
	"hydra/common"
)

func HandleData(message []byte, client *common.Client) {
	response := &common.Response{}
	response.Data.Code = uint32(de.Success)
	//handle message according to yourself
	godog.Debug("[HandleData] receive message:%s", string(message))

	defer func() {
		respByte, err := json.Marshal(response)
		if err != nil {
			godog.Error("[HandleData] response json Marshal occur error: %s", err)
		}

		godog.Debug("[HandleData] deal conn. response: %v", *response)

		client.Socket.WriteMessage(websocket.TextMessage, respByte)
	}()

	data := &common.HeartBeatReq{}
	if err := json.Unmarshal(message, data); err != nil {
		response.Data.Code = uint32(de.SystemError)
		godog.Error("[HandleData] json unmarshal occur error: %s", err)
	}

	response.Id = data.Id
	response.Cmd = data.Cmd

	godog.Debug("[HandleData] deal conn. id: %s, cmd:%s", data.Id, data.Cmd)

	switch data.Cmd {
	case "login":
		loginData := &common.LoginReq{}
		if err := json.Unmarshal(message, loginData); err != nil {
			response.Data.Code = uint32(de.SystemError)
			godog.Error("[HandleData] loginData json unmarshal occur error: %s", err)
			break
		}

		godog.Debug("[HandleData] login uuid: %d", loginData.Uuid)

		client.Uuid = loginData.Uuid
		Hub.Register <- client

		if err := cache.SetUuid(loginData.Uuid); err != nil {
			response.Data.Code = uint32(de.SystemError)
			godog.Error("[HandleData] loginData SetUuid occur error: %s, uuid:%d ", err, loginData.Uuid)
			break
		}

	case "heartbeat":
		godog.Debug("[HandleData] heartbeat uuid:", client.Uuid)

		if client.Uuid > 0 {
			if err := cache.SetUuid(client.Uuid); err != nil {
				response.Data.Code = uint32(de.SystemError)
				godog.Error("[HandleData] heartbeat SetUuid occur error: %s, uuid:%d ", err, client.Uuid)
				break
			}

		} else {
			response.Data.Code = uint32(de.SystemError)
			godog.Debug("[HandleData] heartbeat user not login")
		}

	default:
		response.Data.Code = uint32(de.ParameterError)
		godog.Error("[HandleData] method is invalid. method: %s", data.Cmd)
		break
	}

	return
}
