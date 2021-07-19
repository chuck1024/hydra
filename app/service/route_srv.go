/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/chuck1024/gd"
	"hydra/app/domain"
	"time"
)

func Route(local string, id string, uuid uint64, msg string) (string, error) {
	url := "http://" + local + ":" + gd.Config("Server", "httpPort").MustString("9527") + "/v1/route"
	request := &domain.RouteReq{
		Id:   id,
		Uuid: uuid,
		Msg:  msg,
	}

	resp, _, err := gd.NewHttpClient().Post(url).Timeout(3 * time.Second).Send(request).End()
	if err != nil {
		gd.Error("[Route] send to server occur error: %s", err)
		return "", err
	}

	dataByte, _ := json.Marshal(resp.Body)
	response := &domain.RouteResponse{}
	json.Unmarshal(dataByte, response)

	gd.Debug("[Route] seq:%s", response.Result.Seq)
	seq := response.Result.Seq
	return seq, nil
}
