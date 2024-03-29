/**
 * Copyright 2020 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package boot

import (
	"github.com/gdp-org/gd"
	"hydra/app/route"
	"hydra/app/service"
)

func Run() {
	// init gd
	d := gd.Default()

	// init inject
	Inject()

	// route register
	route.Register(d)

	// start ws
	// todo recover
	go service.Start()

	// gd run
	if err := d.Run(); err != nil {
		gd.Crashf("hydra run occur err:%v", err)
		return
	}
}
