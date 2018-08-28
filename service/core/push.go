/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package core

import (
	"hydra/cache"
	"github.com/chuck1024/godog"
	"hydra/common"
	"encoding/json"
	"github.com/chuck1024/godog/utils"
	"hydra/service/route"
)

func Push(uuid uint64,msg string)(string, error){
	localAddr, err := cache.GetUuid(uuid)
	if err != nil {
		if err == cache.KeyNotExist {
			godog.Debug("[Push] uuid is offline.")
			return "", err
		}

		godog.Error("[Push] GetUuid occur error. err:%s", err)
		return "", err
	}

	if localAddr != utils.GetLocalIP() {
		seq, err := route.Route(localAddr, uuid, msg)
		if err != nil {
			godog.Error("[Push] route occur error: %s", err)
			return "", err
		}
		return seq,nil
	}

	data := &common.TransferData{
		Uuid: uuid,
		Msg:msg,
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		godog.Error("[Push] json marshal occur error:%s",err)
		return "",err
	}

	Hub.SendMsg<-dataByte
	return common.BuildSeq(uuid), err
}