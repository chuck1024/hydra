/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/chuck1024/godog"
	"github.com/chuck1024/godog/utils"
	"github.com/chuck1024/hydra/common"
	"github.com/chuck1024/hydra/dao/cache"
	"github.com/chuck1024/hydra/model"
)

func Push(id string, uuid uint64, msg string) (string, error) {
	localAddr, err := cache.GetUuid(uuid)
	if err != nil {
		if err == cache.KeyNotExist {
			godog.Debug("[Push] uuid is offline.")
			return "", err
		}

		godog.Error("[Push] GetUuid occur error. err:%s", err)
		return "", err
	}

	ip := utils.GetLocalIP()
	//ip := utils.GetLocalIP() +":" + strconv.Itoa(godog.AppConfig.BaseConfig.Server.HttpPort)
	if localAddr != ip {
		seq, err := Route(localAddr, id, uuid, msg)
		if err != nil {
			godog.Error("[Push] route occur error: %s", err)
			return "", err
		}
		return seq, nil
	}

	seq := common.BuildSeq(uuid)

	data := &model.TransferData{
		Seq:  seq,
		Uuid: uuid,
		Msg:  msg,
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		godog.Error("[Push] json marshal occur error:%s", err)
		return "", err
	}

	Hub.SendMsg <- dataByte

	if err := cache.SetPush(id); err != nil {
		godog.Error("[Push] cache ser push occur error: %s", err)
		return "", err
	}

	return seq, err
}
