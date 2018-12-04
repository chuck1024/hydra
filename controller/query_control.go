/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	"github.com/chuck1024/godog"
	de "github.com/chuck1024/godog/error"
	"github.com/chuck1024/godog/net/httplib"
	"hydra/dao/cache"
	"hydra/model"
	"net/http"
)

func QueryControl(rsp http.ResponseWriter, req *http.Request) {
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
	request := &model.QueryReq{}
	response := &model.QueryRsp{}

	defer func() {
		if dErr != nil {
			godog.Error("[QueryControl], errorCode: %d, errMsg: %s", dErr.Code(), dErr.Detail())
		}
		rsp.Write(httplib.LogGetResponseInfo(req, dErr, response))
	}()

	err := httplib.GetRequestBody(req, request)
	if err != nil {
		dErr = de.MakeCodeError(de.ParameterError, err)
		return
	}

	godog.Info("[QueryControl] received request: %v", *request)

	_, err = cache.GetUuid(request.Uuid)
	if err != nil {
		if err == cache.KeyNotExist {
			response.IsOnline = false
			godog.Debug("[QueryControl] uuid[%d] is offline.", request.Uuid)
			return
		}
		godog.Error("[QueryControl] cache get uuid occur error:%s", err)
		dErr = de.MakeCodeError(de.SystemError, err)
		return
	}

	response.IsOnline = true
	return
}
