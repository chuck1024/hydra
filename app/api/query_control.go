/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package api

import (
	"github.com/chuck1024/gd"
	de "github.com/chuck1024/gd/derror"
	"github.com/gin-gonic/gin"
	"hydra/app/domain"
	"hydra/app/model"
	"hydra/app/service/sp"
)

func QueryControl(c *gin.Context, req *domain.QueryReq) (code int, message string, err error, ret *domain.QueryRsp) {
	ret = &domain.QueryRsp{}

	_, err = sp.Get().UidCache.GetUuid(req.Uuid)
	if err != nil {
		if err == model.KeyNotExist {
			ret.IsOnline = false
			gd.Debug("[QueryControl] uuid[%d] is offline.", req.Uuid)
			return de.Success, "", nil, ret
		}
		gd.Error("[QueryControl] model get uuid occur error:%s", err)
		return de.SystemError, err.Error(), err, ret
	}

	ret.IsOnline = true
	return de.Success, "", nil, ret
}
