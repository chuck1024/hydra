/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package api

import (
	"errors"
	"github.com/gdp-org/gd"
	de "github.com/gdp-org/gd/derror"
	"github.com/gin-gonic/gin"
	"hydra/app/domain"
	"hydra/app/model"
	"hydra/app/service"
	"hydra/app/service/sp"
)

func PushControl(c *gin.Context, req *domain.PushReq) (code int, message string, err error, ret *domain.PushRsp) {
	ret = &domain.PushRsp{}

	if sp.Get().UidCache.GetPush(req.Id) {
		gd.Error("[Push] model get push, id[%s] is exist", req.Id)
		return de.ParameterError, "already sent id", errors.New("already sent id"), ret
	}

	seq, err := service.Push(req.Id, req.Uuid, req.Msg)
	if err != nil {
		if err == model.KeyNotExist {
			gd.Debug("[PushControl] uuid[%d] is offline.", req.Uuid)
			return domain.Offline, err.Error(), err, ret
		}
		gd.Error("[PushControl] push occur error:%s", err)
		return de.SystemError, err.Error(), err, ret
	}

	ret.Seq = seq
	return de.Success, "", nil, ret
}
