/**
 * Copyright 2020 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package route

import (
	"github.com/chuck1024/gd"
	"github.com/chuck1024/gd/net/dhttp"
	"github.com/gin-gonic/gin"
	"hydra/app/api"
	"sync"
)

var (
	initOnce sync.Once
)

func Register(e *gd.Engine) {
	e.HttpServer.SetInit(func(g *gin.Engine) error {
		r := g.Group("")

		r.Use(
			dhttp.GlFilter(),
			dhttp.StatFilter(),
			dhttp.GroupFilter(),
			dhttp.Logger(gd.Config("Server", "serverName").String()),
		)
		return route(e, r)
	})
}

func route(e *gd.Engine, r *gin.RouterGroup) error {
	var ret error
	initOnce.Do(func() {
		g := r.Group("v1")

		e.HttpServer.POST(g, "push", api.PushControl)
		e.HttpServer.POST(g, "query", api.QueryControl)
		e.HttpServer.POST(g, "route", api.RouteControl)

		if ret = e.HttpServer.CheckHandle(); ret != nil {
			return
		}

		g.GET("hydra", api.WsControl)
	})
	return ret
}
