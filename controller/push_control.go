/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	"errors"
	"github.com/chuck1024/godog"
	de "github.com/chuck1024/godog/error"
	"github.com/chuck1024/godog/net/httplib"
	"hydra/common"
	"hydra/dao/cache"
	"hydra/model"
	"hydra/service"
	"net/http"
)

func PushControl(rsp http.ResponseWriter, req *http.Request) {
	rsp.Header().Add("Access-Control-Allow-Origin", httplib.CONTENT_ALL)
	rsp.Header().Add("Content-Type", httplib.CONTENT_JSON)

	if req.Method == http.MethodOptions {
		rsp.WriteHeader(http.StatusOK)
		return
	} else if req.Method != http.MethodPost {
		rsp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var dErr *de.CodeError
	request := &model.PushReq{}
	response := &model.PushRsp{}

	defer func() {
		if dErr != nil {
			godog.Error("[PushControl], errorCode: %d, errMsg: %s", dErr.Code(), dErr.Detail())
		}
		rsp.Write(httplib.LogGetResponseInfo(req, dErr, response))
	}()

	err := httplib.GetRequestBody(req, &request)
	if err != nil {
		dErr = de.MakeCodeError(de.ParameterError, err)
		return
	}

	godog.Info("[PushControl] received request: %v", *request)

	if cache.GetPush(request.Id) {
		godog.Error("[Push] cache get push, id[%s] is exist", request.Id)
		dErr = de.MakeCodeError(de.ParameterError, errors.New("already sent id"))
		return
	}

	seq, err := service.Push(request.Id, request.Uuid, request.Msg)
	if err != nil {
		if err == cache.KeyNotExist {
			godog.Debug("[PushControl] uuid[%d] is offline.", request.Uuid)
			dErr = de.MakeCodeError(common.Offline, err)
			return
		}
		godog.Error("[PushControl] push occur error:%s", err)
		dErr = de.MakeCodeError(de.SystemError, err)
		return
	}

	response.Seq = seq
	return
}
