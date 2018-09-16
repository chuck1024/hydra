/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package main

import (
	"github.com/chuck1024/godog"
	"github.com/chuck1024/godog/dao/cache"
	"hydra/controller"
	"hydra/model/service"
)

func register() {
	godog.AppHttp.AddHttpHandler("/hydra", controller.WsControl)
	godog.AppHttp.AddHttpHandler("/push", controller.PushControl)
	godog.AppHttp.AddHttpHandler("/query", controller.QueryControl)
	godog.AppHttp.AddHttpHandler("/route", controller.RouteControl)
}

func main() {
	url, _ := godog.AppConfig.String("redis")
	cache.Init(url)

	register()

	go service.Start()

	if err := godog.Run(); err != nil {
		godog.Error("Error occurs, error = %s", err.Error())
		return
	}
}
