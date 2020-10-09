/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package api

import (
	de "github.com/chuck1024/gd/derror"
	"github.com/gin-gonic/gin"
	"hydra/app/service"
	"hydra/libray"
)

func RouteControl(c *gin.Context, req *libray.RouteReq) (code int, message string, err error, ret *libray.RouteRsp) {
	ret = &libray.RouteRsp{}

	seq, err := service.Push(req.Id, req.Uuid, req.Msg)
	if err != nil {
		return de.SystemError, err.Error(), err, ret
	}

	ret.Seq = seq
	return de.Success, "", nil, ret
}
