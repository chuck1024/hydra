/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/chuck1024/doglog"
	"github.com/chuck1024/godog/utils"
	"github.com/chuck1024/hydra/common"
	"github.com/chuck1024/hydra/dao/cache"
	"github.com/chuck1024/hydra/model"
)

func Push(id string, uuid uint64, msg string) (string, error) {
	localAddr, err := cache.GetUuid(uuid)
	if err != nil {
		if err == cache.KeyNotExist {
			doglog.Debug("[Push] uuid is offline.")
			return "", err
		}

		doglog.Error("[Push] GetUuid occur error. err:%s", err)
		return "", err
	}

	ip := utils.GetLocalIP()
	//ip := utils.GetLocalIP() +":" + strconv.Itoa(dog.Config.BaseConfig.Server.HttpPort)
	if localAddr != ip {
		seq, err := Route(localAddr, id, uuid, msg)
		if err != nil {
			doglog.Error("[Push] route occur error: %s", err)
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
		doglog.Error("[Push] json marshal occur error:%s", err)
		return "", err
	}

	Hub.SendMsg <- dataByte

	if err := cache.SetPush(id); err != nil {
		doglog.Error("[Push] cache ser push occur error: %s", err)
		return "", err
	}

	return seq, err
}
