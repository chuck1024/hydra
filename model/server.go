/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package model

// push msg
type PushReq struct {
	Id   string
	Uuid uint64
	Msg  string
}

type PushRsp struct {
	Seq string
}

// query isOnline
type QueryReq struct {
	Uuid uint64
}

type QueryRsp struct {
	IsOnline bool
}
