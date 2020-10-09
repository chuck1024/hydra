/**
 * Copyright 2020 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package boot

import (
	"github.com/chuck1024/gd"
	"github.com/chuck1024/gd/dlog"
	"github.com/chuck1024/gd/runtime/inject"
	"hydra/route"
)

func Run() {
	// init gd
	d := gd.Default()

	// init inject
	inject.InitDefault()
	inject.SetLogger(dlog.Global)
	defer inject.Close()
	Inject(d)

	// route register
	route.Register(d)

	// gd run
	if err := d.Run(); err != nil {
		dlog.Crashf("hydra run occur err:%v", err)
		return
	}
}
