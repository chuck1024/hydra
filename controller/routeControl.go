/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package controller

import (
	"github.com/chuck1024/godog"
	de "github.com/chuck1024/godog/error"
	"github.com/chuck1024/godog/net/httplib"
	"hydra/common"
	"hydra/service/core"
	"net/http"
)

func RouteControl(rsp http.ResponseWriter, req *http.Request) {
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
	request := &common.RouteReq{}
	response := &common.RouteRsp{}

	defer func() {
		if dErr != nil {
			godog.Error("[RouteControl], errorCode: %d, errMsg: %s", dErr.Code(), dErr.Detail())
		}
		rsp.Write(httplib.LogGetResponseInfo(req, dErr, response))
	}()

	err := httplib.GetRequestBody(req, &request)
	if err != nil {
		dErr = de.MakeCodeError(de.ParameterError, err)
		return
	}

	godog.Info("[RouteControl] received request: %#v", request)

	seq, err := core.Push(request.Id, request.Uuid, request.Msg)
	if err != nil {
		dErr = de.MakeCodeError(de.SystemError, err)
		return
	}

	response.Seq = seq
}
