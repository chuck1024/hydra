/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package model

////////// to client /////////////
// login
type LoginReq struct {
	Id   string `json:"id"`
	Cmd  string `json:"cmd"` // login
	Uuid uint64 `json:"uuid"`
}

// heartbeat
type HeartBeatReq struct {
	Id  string `json:"id"`
	Cmd string `json:"cmd"` // heartbeat
}

// push
type PushClientReq struct {
	Id  string `json:"id"`
	Cmd string `json:"cmd"`
	Msg string `json:"msg"`
}

type Response struct {
	Id   string `json:"id"`
	Cmd  string `json:"cmd"`
	Data struct {
		Code   uint32 `json:"code"`
		Result string `json:"result"`
	}
}
