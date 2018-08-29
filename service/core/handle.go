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
		if response.Cmd == "push" {
			return
		}

		response.Data.Result = de.GetErrorType(int(response.Data.Code))
		respByte, err := json.Marshal(response)
		if err != nil {
			godog.Error("[HandleData] response json Marshal occur error: %s", err)
		}

		godog.Debug("[HandleData] handle data. response: %v", *response)

		client.Socket.WriteMessage(websocket.TextMessage, respByte)
	}()

	data := &common.HeartBeatReq{}
	if err := json.Unmarshal(message, data); err != nil {
		response.Data.Code = uint32(de.SystemError)
		godog.Error("[HandleData] json unmarshal occur error: %s", err)
	}

	response.Id = data.Id
	response.Cmd = data.Cmd

	godog.Debug("[HandleData] handle data. id: %s, cmd:%s", data.Id, data.Cmd)

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
		godog.Debug("[HandleData] heartbeat uuid: %d", client.Uuid)

		if client.Uuid > 0 {
			if _, err := cache.GetUuid(client.Uuid); err != nil {
				if err == cache.KeyNotExist {
					response.Data.Code = uint32(de.SystemError)
					godog.Debug("[HandleData] heartbeat user not login")
					break
				}

				response.Data.Code = uint32(de.SystemError)
				godog.Error("[HandleData] heartbeat SetUuid occur error: %s, uuid:%d ", err, client.Uuid)
				break
			}

			if err := cache.SetUuid(client.Uuid); err != nil {
				response.Data.Code = uint32(de.SystemError)
				godog.Error("[HandleData] heartbeat SetUuid occur error: %s, uuid:%d ", err, client.Uuid)
				break
			}

		} else {
			response.Data.Code = uint32(de.SystemError)
			godog.Debug("[HandleData] heartbeat user not login")
		}

	case "push":
		rsp := &common.Response{}
		if err := json.Unmarshal(message, rsp); err != nil {
			response.Data.Code = uint32(de.SystemError)
			godog.Error("[HandleData] push response json unmarshal occur error: %s", err)
			break
		}

		godog.Debug("[HandleData] push response: %v", rsp)

	default:
		response.Data.Code = uint32(de.ParameterError)
		godog.Error("[HandleData] method is invalid. method: %s", data.Cmd)
		break
	}

	return
}
