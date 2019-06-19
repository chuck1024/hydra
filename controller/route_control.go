/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	de "github.com/chuck1024/godog/error"
	"github.com/chuck1024/hydra/model"
	"github.com/chuck1024/hydra/service"
	"github.com/gin-gonic/gin"
)

func RouteControl(c *gin.Context, req *model.RouteReq) (code int, message string, err error, ret *model.RouteRsp) {
	ret = &model.RouteRsp{}

	seq, err := service.Push(req.Id, req.Uuid, req.Msg)
	if err != nil {
		return de.SystemError, err.Error(), err, ret
	}

	ret.Seq = seq
	return de.Success, "", nil, ret
}
