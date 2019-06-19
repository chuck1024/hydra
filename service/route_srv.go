/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/chuck1024/doglog"
	"github.com/chuck1024/godog"
	"github.com/chuck1024/hydra/model"
	"strconv"
	"time"
)

func Route(local string, id string, uuid uint64, msg string) (string, error) {
	dog := godog.Default()
	url := "http://" + local + ":" + strconv.Itoa(dog.Config.BaseConfig.Server.HttpPort)
	//url := "http://" + local  + "/route"
	request := &model.RouteReq{
		Id:   id,
		Uuid: uuid,
		Msg:  msg,
	}

	client := dog.NewHttpClient(time.Duration(0), url)
	resp, _, err := client.Method("POST", "route", nil, request)
	if err != nil {
		doglog.Error("[Route] send to server occur error: %s", err)
		return "", err
	}

	dataByte, _ := json.Marshal(resp.Body)
	response := &model.RouteRsp{}
	json.Unmarshal(dataByte, response)

	doglog.Debug("[Route] seq:%s", response.Seq)
	seq := response.Seq
	return seq, nil
}
