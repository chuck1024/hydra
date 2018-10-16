/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package model

// push transfer
type TransferData struct {
	Seq  string `json:"seq"`
	Uuid uint64 `json:"uuid"`
	Msg  string `json:"msg"`
}
