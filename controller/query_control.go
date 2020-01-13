/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	"github.com/chuck1024/doglog"
	de "github.com/chuck1024/godog/error"
	"github.com/gin-gonic/gin"
	"hydra/common"
	"hydra/dao/cache"
)

func QueryControl(c *gin.Context, req *common.QueryReq) (code int, message string, err error, ret *common.QueryRsp) {
	ret = &common.QueryRsp{}

	_, err = cache.GetUuid(req.Uuid)
	if err != nil {
		if err == cache.KeyNotExist {
			ret.IsOnline = false
			doglog.Debug("[QueryControl] uuid[%d] is offline.", req.Uuid)
			return de.Success, "", nil, ret
		}
		doglog.Error("[QueryControl] cache get uuid occur error:%s", err)
		return de.SystemError, err.Error(), err, ret
	}

	ret.IsOnline = true
	return de.Success, "", nil, ret
}
