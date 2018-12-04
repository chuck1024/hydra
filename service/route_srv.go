/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/chuck1024/godog"
	"github.com/chuck1024/godog/net/httplib"
	"hydra/model"
	"strconv"
)

func Route(local string, id string, uuid uint64, msg string) (string, error) {
	url := "http://" + local + ":" + strconv.Itoa(godog.AppConfig.BaseConfig.Server.HttpPort) + "/route"
	//url := "http://" + local  + "/route"
	request := &model.RouteReq{
		Id:   id,
		Uuid: uuid,
		Msg:  msg,
	}

	resp := &httplib.ResponseData{}
	err := httplib.SendToServer(httplib.HttpPost, url, nil, nil, request, resp)
	if err != nil {
		godog.Error("[Route] send to server occur error: %s", err)
		return "", err
	}

	dataByte, _ := json.Marshal(resp.Data)
	response := &model.RouteRsp{}
	json.Unmarshal(dataByte, response)

	godog.Debug("[Route] seq:%s", response.Seq)
	seq := response.Seq
	return seq, nil
}
