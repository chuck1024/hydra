/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	de "github.com/chuck1024/godog/error"
	"github.com/chuck1024/godog/net/httplib"
	"net/http"
	"hydra/common"
	"github.com/chuck1024/godog"
	"hydra/service/core"
)

func RouteControl(rsp http.ResponseWriter, req *http.Request){
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
	request := &common.RouteReq{}
	response := &common.RouteRsp{}

	defer func() {
		if dErr != nil {
			godog.Error("[RouteControl], errorCode: %d, errMsg: %s", dErr.Code(), dErr.Detail())
		}
		rsp.Write(httplib.LogGetResponseInfo(req, dErr, rsp))
	}()

	err := httplib.GetRequestBody(req, &request)
	if err != nil {
		dErr = de.MakeCodeError(de.ParameterError, err)
		return
	}

	godog.Info("[RouteControl] received request: %#v", request)

	seq, err := core.Push(request.Uuid,request.Msg)
	if err != nil {
		dErr = de.MakeCodeError(de.SystemError, err)
		return
	}

	response.Seq = seq
}