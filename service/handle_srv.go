/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/chuck1024/doglog"
	de "github.com/chuck1024/godog/error"
	"github.com/gorilla/websocket"
	"hydra/common"
	"hydra/dao/cache"
	"hydra/model"
)

func HandleData(message []byte, client *model.Client) {
	response := &model.Response{}
	response.Data.Code = uint32(de.Success)
	//handle message according to yourself
	doglog.Debug("[HandleData] receive message:%s", string(message))

	defer func() {
		if response.Cmd == common.PushCmd {
			return
		}

		response.Data.Result = de.GetErrorType(int(response.Data.Code))
		respByte, err := json.Marshal(response)
		if err != nil {
			doglog.Error("[HandleData] response json Marshal occur error: %s", err)
		}

		doglog.Debug("[HandleData] handle data. response: %v", *response)

		client.Socket.WriteMessage(websocket.TextMessage, respByte)
	}()

	data := &model.HeartBeatReq{}
	if err := json.Unmarshal(message, data); err != nil {
		response.Data.Code = uint32(de.SystemError)
		doglog.Error("[HandleData] json unmarshal occur error: %s", err)
	}

	response.Id = data.Id
	response.Cmd = data.Cmd

	doglog.Debug("[HandleData] handle data. id: %s, cmd:%s", data.Id, data.Cmd)

	switch data.Cmd {
	case common.LoginCmd:
		loginData := &model.LoginReq{}
		if err := json.Unmarshal(message, loginData); err != nil {
			response.Data.Code = uint32(de.SystemError)
			doglog.Error("[HandleData] loginData json unmarshal occur error: %s", err)
			break
		}

		doglog.Debug("[HandleData] login uuid: %d", loginData.Uuid)

		client.Uuid = loginData.Uuid
		Hub.Register <- client

		if err := cache.SetUuid(loginData.Uuid); err != nil {
			response.Data.Code = uint32(de.SystemError)
			doglog.Error("[HandleData] loginData SetUuid occur error: %s, uuid:%d ", err, loginData.Uuid)
			break
		}

	case common.HeartbeatCmd:
		doglog.Debug("[HandleData] heartbeat uuid: %d", client.Uuid)

		if client.Uuid > 0 {
			if _, err := cache.GetUuid(client.Uuid); err != nil {
				if err == cache.KeyNotExist {
					response.Data.Code = uint32(de.SystemError)
					doglog.Debug("[HandleData] heartbeat user not login")
					break
				}

				response.Data.Code = uint32(de.SystemError)
				doglog.Error("[HandleData] heartbeat SetUuid occur error: %s, uuid:%d ", err, client.Uuid)
				break
			}

			if err := cache.SetUuid(client.Uuid); err != nil {
				response.Data.Code = uint32(de.SystemError)
				doglog.Error("[HandleData] heartbeat SetUuid occur error: %s, uuid:%d ", err, client.Uuid)
				break
			}

		} else {
			response.Data.Code = uint32(de.SystemError)
			doglog.Debug("[HandleData] heartbeat user not login")
		}

	case common.PushCmd:
		rsp := &model.Response{}
		if err := json.Unmarshal(message, rsp); err != nil {
			response.Data.Code = uint32(de.SystemError)
			doglog.Error("[HandleData] push response json unmarshal occur error: %s", err)
			break
		}

		doglog.Debug("[HandleData] push response: %v", rsp)

	default:
		response.Data.Code = uint32(de.ParameterError)
		doglog.Error("[HandleData] method is invalid. method: %s", data.Cmd)
		break
	}

	return
}
