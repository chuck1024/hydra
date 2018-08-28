/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package main

import (
	"github.com/chuck1024/godog"
	"hydra/controller"
	"hydra/service/core"
)

func register() {
	godog.AppHttp.AddHttpHandler("/hydra", controller.WsControl)
	godog.AppHttp.AddHttpHandler("/push", controller.PushControl)
	godog.AppHttp.AddHttpHandler("/query", controller.QueryControl)
	godog.AppHttp.AddHttpHandler("/route", controller.RouteControl)
}

func main() {
	register()

	go core.Start()

	if err := godog.Run(); err != nil {
		godog.Error("Error occurs, error = %s", err.Error())
		return
	}
}
