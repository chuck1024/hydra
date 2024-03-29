/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/gdp-org/gd"
	de "github.com/gdp-org/gd/derror"
	"github.com/gdp-org/gd/runtime/gl"
	"github.com/gorilla/websocket"
	"hydra/app/domain"
	"hydra/app/model"
	"hydra/app/service/sp"
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

	response := &domain.Response{}
	response.Data.Code = uint32(de.Success)
	//handle message according to yourself
	gd.Debug("[HandleData] receive message:%s", string(message))

	defer func() {
		if response.Cmd == domain.PushCmd {
			return
		}

		response.Data.Result = de.GetErrorType(int(response.Data.Code))
		respByte, err := json.Marshal(response)
		if err != nil {
			gd.Error("[HandleData] response json Marshal occur error: %s", err)
		}

		gd.Debug("[HandleData] handle data. response: %v", *response)

		client.Socket.WriteMessage(websocket.TextMessage, respByte)
	}()

	data := &domain.HeartBeatReq{}
	if err := json.Unmarshal(message, data); err != nil {
		response.Data.Code = uint32(de.SystemError)
		gd.Error("[HandleData] json unmarshal occur error: %s", err)
	}

	response.Id = data.Id
	response.Cmd = data.Cmd

	gd.Debug("[HandleData] handle data. id: %s, cmd:%s", data.Id, data.Cmd)

	switch data.Cmd {
	case domain.LoginCmd:
		loginData := &domain.LoginReq{}
		if err := json.Unmarshal(message, loginData); err != nil {
			response.Data.Code = uint32(de.SystemError)
			gd.Error("[HandleData] loginData json unmarshal occur error: %s", err)
			break
		}

		gd.Debug("[HandleData] login uuid: %d", loginData.Uuid)

		client.Uuid = loginData.Uuid
		Hub.Register <- client

		if err := sp.Get().UidCache.SetUuid(loginData.Uuid); err != nil {
			response.Data.Code = uint32(de.SystemError)
			gd.Error("[HandleData] loginData SetUuid occur error: %s, uuid:%d ", err, loginData.Uuid)
			break
		}

	case domain.HeartbeatCmd:
		gd.Debug("[HandleData] heartbeat uuid: %d", client.Uuid)

		if client.Uuid > 0 {
			if _, err := sp.Get().UidCache.GetUuid(client.Uuid); err != nil {
				if err == model.KeyNotExist {
					response.Data.Code = uint32(de.SystemError)
					gd.Debug("[HandleData] heartbeat user not login")
					break
				}

				response.Data.Code = uint32(de.SystemError)
				gd.Error("[HandleData] heartbeat SetUuid occur error: %s, uuid:%d ", err, client.Uuid)
				break
			}

			if err := sp.Get().UidCache.SetUuid(client.Uuid); err != nil {
				response.Data.Code = uint32(de.SystemError)
				gd.Error("[HandleData] heartbeat SetUuid occur error: %s, uuid:%d ", err, client.Uuid)
				break
			}

		} else {
			response.Data.Code = uint32(de.SystemError)
			gd.Debug("[HandleData] heartbeat user not login")
		}

	case domain.PushCmd:
		rsp := &domain.Response{}
		if err := json.Unmarshal(message, rsp); err != nil {
			response.Data.Code = uint32(de.SystemError)
			gd.Error("[HandleData] push response json unmarshal occur error: %s", err)
			break
		}

		gd.Debug("[HandleData] push response: %v", rsp)

	default:
		response.Data.Code = uint32(de.ParameterError)
		gd.Error("[HandleData] method is invalid. method: %s", data.Cmd)
		break
	}

	return
}
