/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package route

import (
	"github.com/chuck1024/godog"
	"github.com/chuck1024/godog/net/httplib"
	"hydra/common"
	"strconv"
)

func Route(local string, id string, uuid uint64, msg string) (string, error) {
	url := local + strconv.Itoa(godog.AppConfig.BaseConfig.Server.HttpPort) + "/route"
	request := &common.RouteReq{
		Id:   id,
		Uuid: uuid,
		Msg:  msg,
	}
	response := &common.RouteRsp{}
	err := httplib.SendToServer(httplib.HttpPost, url, nil, nil, request, response)
	if err != nil {
		godog.Error("[Route] send to server occur error: %s", err)
		return "", err
	}

	seq := response.Seq
	return seq, nil
}
