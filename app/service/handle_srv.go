/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	de "github.com/chuck1024/gd/derror"
	"github.com/chuck1024/gd/dlog"
	"github.com/chuck1024/gd/runtime/gl"
	"github.com/gorilla/websocket"
	"hydra/app/model"
	"hydra/app/service/sp"
	"hydra/libray"
	"strconv"
	"time"
)

func HandleData(message []byte, client *Client) {
	gl.Init()
	defer gl.Close()
	st := time.Now()
	traceId := strconv.FormatInt(st.UnixNano(), 10)
	gl.Set(gl.LogId, traceId)
	gl.Set(gl.ClientIp, client.Socket.RemoteAddr().String())

	response := &libray.Response{}
	response.Data.Code = uint32(de.Success)
	//handle message according to yourself
	dlog.Debug("[HandleData] receive message:%s", string(message))

	defer func() {
		if response.Cmd == libray.PushCmd {
			return
		}

		response.Data.Result = de.GetErrorType(int(response.Data.Code))
		respByte, err := json.Marshal(response)
		if err != nil {
			dlog.Error("[HandleData] response json Marshal occur error: %s", err)
		}

		dlog.Debug("[HandleData] handle data. response: %v", *response)

		client.Socket.WriteMessage(websocket.TextMessage, respByte)
	}()

	data := &libray.HeartBeatReq{}
	if err := json.Unmarshal(message, data); err != nil {
		response.Data.Code = uint32(de.SystemError)
		dlog.Error("[HandleData] json unmarshal occur error: %s", err)
	}

	response.Id = data.Id
	response.Cmd = data.Cmd

	dlog.Debug("[HandleData] handle data. id: %s, cmd:%s", data.Id, data.Cmd)

	switch data.Cmd {
	case libray.LoginCmd:
		loginData := &libray.LoginReq{}
		if err := json.Unmarshal(message, loginData); err != nil {
			response.Data.Code = uint32(de.SystemError)
			dlog.Error("[HandleData] loginData json unmarshal occur error: %s", err)
			break
		}

		dlog.Debug("[HandleData] login uuid: %d", loginData.Uuid)

		client.Uuid = loginData.Uuid
		Hub.Register <- client

		if err := sp.Get().UidCache.SetUuid(loginData.Uuid); err != nil {
			response.Data.Code = uint32(de.SystemError)
			dlog.Error("[HandleData] loginData SetUuid occur error: %s, uuid:%d ", err, loginData.Uuid)
			break
		}

	case libray.HeartbeatCmd:
		dlog.Debug("[HandleData] heartbeat uuid: %d", client.Uuid)

		if client.Uuid > 0 {
			if _, err := sp.Get().UidCache.GetUuid(client.Uuid); err != nil {
				if err == model.KeyNotExist {
					response.Data.Code = uint32(de.SystemError)
					dlog.Debug("[HandleData] heartbeat user not login")
					break
				}

				response.Data.Code = uint32(de.SystemError)
				dlog.Error("[HandleData] heartbeat SetUuid occur error: %s, uuid:%d ", err, client.Uuid)
				break
			}

			if err := sp.Get().UidCache.SetUuid(client.Uuid); err != nil {
				response.Data.Code = uint32(de.SystemError)
				dlog.Error("[HandleData] heartbeat SetUuid occur error: %s, uuid:%d ", err, client.Uuid)
				break
			}

		} else {
			response.Data.Code = uint32(de.SystemError)
			dlog.Debug("[HandleData] heartbeat user not login")
		}

	case libray.PushCmd:
		rsp := &libray.Response{}
		if err := json.Unmarshal(message, rsp); err != nil {
			response.Data.Code = uint32(de.SystemError)
			dlog.Error("[HandleData] push response json unmarshal occur error: %s", err)
			break
		}

		dlog.Debug("[HandleData] push response: %v", rsp)

	default:
		response.Data.Code = uint32(de.ParameterError)
		dlog.Error("[HandleData] method is invalid. method: %s", data.Cmd)
		break
	}

	return
}
