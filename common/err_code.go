/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package common

import (
	"github.com/chuck1024/godog/error"
)

var (
	Offline    = 601
	OfflineStr = "offline"
)

func init() {
	error.ErrMap[Offline] = OfflineStr
}
