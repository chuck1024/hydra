/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package service

import (
	"encoding/json"
	"github.com/chuck1024/gd/dlog"
	"github.com/chuck1024/gd/utls/network"
	"hydra/app/model"
	"hydra/app/service/sp"
	"hydra/libray"
)

func Push(id string, uuid uint64, msg string) (string, error) {
	localAddr, err := sp.Get().UidCache.GetUuid(uuid)
	if err != nil {
		if err == model.KeyNotExist {
			dlog.Debug("[Push] uuid is offline.")
			return "", err
		}

		dlog.Error("[Push] GetUuid occur error. err:%s", err)
		return "", err
	}

	ip := network.GetLocalIP()
	//ip := utils.GetLocalIP() +":" + strconv.Itoa(dog.Config.BaseConfig.Server.HttpPort)
	if localAddr != ip {
		seq, err := Route(localAddr, id, uuid, msg)
		if err != nil {
			dlog.Error("[Push] route occur error: %s", err)
			return "", err
		}
		return seq, nil
	}

	seq := libray.BuildSeq(uuid)

	data := &libray.TransferData{
		Seq:  seq,
		Uuid: uuid,
		Msg:  msg,
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		dlog.Error("[Push] json marshal occur error:%s", err)
		return "", err
	}

	Hub.SendMsg <- dataByte

	if err := sp.Get().UidCache.SetPush(id); err != nil {
		dlog.Error("[Push] model ser push occur error: %s", err)
		return "", err
	}

	return seq, err
}
