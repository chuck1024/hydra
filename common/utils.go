/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package common

import (
	"fmt"
	"time"
)

func BuildSeq(uuid uint64) string {
	return fmt.Sprintf("%s%d%d",  time.Now().Format("20060102150405"), uuid, 1)
}
