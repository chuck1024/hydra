/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package domain

import (
	"github.com/gdp-org/gd/derror"
)

var (
	Offline    = 601
	OfflineStr = "offline"
)

func init() {
	derror.ErrMap[Offline] = OfflineStr
}
