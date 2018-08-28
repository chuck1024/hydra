/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	de "github.com/chuck1024/godog/error"
	"github.com/chuck1024/godog/net/httplib"
	"github.com/chuck1024/godog"
	"net/http"
	"hydra/common"
	"hydra/cache"
)

func QueryControl(rsp http.ResponseWriter, req *http.Request){
	rsp.Header().Add("Access-Control-Allow-Origin", httplib.CONTENT_ALL)
	rsp.Header().Add("Content-Type", httplib.CONTENT_JSON)

	if req.Method == http.MethodOptions {
		rsp.WriteHeader(http.StatusOK)
		return
	} else if req.Method != http.MethodPost {
		rsp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dErr := &de.CodeError{}
	request := &common.QueryReq{}
	response := &common.QueryRsp{}

	defer func() {
		if dErr != nil {
			godog.Error("[QueryControl], errorCode: %d, errMsg: %s", dErr.Code(), dErr.Detail())
		}
		rsp.Write(httplib.LogGetResponseInfo(req, dErr, rsp))
	}()

	err := httplib.GetRequestBody(req, &request)
	if err != nil {
		dErr = de.MakeCodeError(de.ParameterError, err)
		return
	}

	godog.Info("[QueryControl] received request: %#v", request)

	_, err = cache.GetUuid(request.Uuid)
	if err != nil {
		if err == cache.KeyNotExist {
			response.IsOnline = false
			godog.Debug("[QueryControl] uuid[%d] is offline.", request.Uuid)
			return
		}
		godog.Error("[QueryControl] cache get uuid occur error:%s",err)
		dErr = de.MakeCodeError(de.SystemError,err)
		return
	}

	response.IsOnline = true
	return
}