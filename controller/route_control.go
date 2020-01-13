/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	de "github.com/chuck1024/godog/error"
	"github.com/gin-gonic/gin"
	"hydra/common"
	"hydra/service"
)

func RouteControl(c *gin.Context, req *common.RouteReq) (code int, message string, err error, ret *common.RouteRsp) {
	ret = &common.RouteRsp{}

	seq, err := service.Push(req.Id, req.Uuid, req.Msg)
	if err != nil {
		return de.SystemError, err.Error(), err, ret
	}

	ret.Seq = seq
	return de.Success, "", nil, ret
}
