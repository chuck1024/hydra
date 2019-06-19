/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	"errors"
	"github.com/chuck1024/doglog"
	de "github.com/chuck1024/godog/error"
	"github.com/chuck1024/hydra/common"
	"github.com/chuck1024/hydra/dao/cache"
	"github.com/chuck1024/hydra/model"
	"github.com/chuck1024/hydra/service"
	"github.com/gin-gonic/gin"
)

func PushControl(c *gin.Context, req *model.PushReq) (code int, message string, err error, ret *model.PushRsp) {
	ret = &model.PushRsp{}

	if cache.GetPush(req.Id) {
		doglog.Error("[Push] cache get push, id[%s] is exist", req.Id)
		return de.ParameterError, "already sent id", errors.New("already sent id"), ret
	}

	seq, err := service.Push(req.Id, req.Uuid, req.Msg)
	if err != nil {
		if err == cache.KeyNotExist {
			doglog.Debug("[PushControl] uuid[%d] is offline.", req.Uuid)
			return common.Offline, err.Error(), err, ret
		}
		doglog.Error("[PushControl] push occur error:%s", err)
		return de.SystemError, err.Error(), err, ret
	}

	ret.Seq = seq
	return de.Success, "", nil, ret
}
