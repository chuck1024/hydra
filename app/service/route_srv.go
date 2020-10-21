/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/chuck1024/gd"
	"github.com/chuck1024/gd/config"
	"hydra/libray"
	"time"
)

func Route(local string, id string, uuid uint64, msg string) (string, error) {
	url := "http://" + local + ":" + config.Config().Section("Server").Key("httpPort").MustString("9527")
	//url := "http://" + local  + "/route"
	request := &libray.RouteReq{
		Id:   id,
		Uuid: uuid,
		Msg:  msg,
	}

	client := gd.NewHttpClient(time.Duration(0), url)
	resp, _, err := client.Method("POST", "route", nil, request)
	if err != nil {
		gd.Error("[Route] send to server occur error: %s", err)
		return "", err
	}

	dataByte, _ := json.Marshal(resp.Body)
	response := &libray.RouteRsp{}
	json.Unmarshal(dataByte, response)

	gd.Debug("[Route] seq:%s", response.Seq)
	seq := response.Seq
	return seq, nil
}
