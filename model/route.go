/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package model

// route
type RouteReq struct {
	Id   string
	Uuid uint64
	Msg  string
}

type RouteRsp struct {
	Seq string
}
