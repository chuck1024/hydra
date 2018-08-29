/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package route

import (
	"encoding/json"
	"github.com/chuck1024/godog"
	"github.com/chuck1024/godog/net/httplib"
	"hydra/common"
	"strconv"
)

func Route(local string, id string, uuid uint64, msg string) (string, error) {
	url := "http://" + local + ":" + strconv.Itoa(godog.AppConfig.BaseConfig.Server.HttpPort) + "/route"
	//url := "http://" + local  + "/route"
	request := &common.RouteReq{
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
	response := &common.RouteRsp{}
	json.Unmarshal(dataByte, response)

	godog.Debug("[Route] seq:%s", response.Seq)
	seq := response.Seq
	return seq, nil
}
