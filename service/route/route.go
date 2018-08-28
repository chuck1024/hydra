/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package route

import (
	"github.com/chuck1024/godog/net/httplib"
	"hydra/common"
	"github.com/chuck1024/godog"
	"strconv"
)

func Route(local string,uuid uint64, msg string)(string, error){
	url := local + strconv.Itoa(godog.AppConfig.BaseConfig.Server.HttpPort) + "/route"
	request := &common.RouteReq{
		Uuid:uuid,
		Msg:msg,
	}
	response := &common.RouteRsp{}
	err := httplib.SendToServer(httplib.HttpPost,url,nil,nil,request,response)
	if err != nil{
		godog.Error("[Route] send to server occur error: %s",err)
		return "", err
	}

	seq := response.Seq
	return seq, nil
}