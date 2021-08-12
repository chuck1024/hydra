/**
 * Copyright 2020 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package route

import (
	"github.com/gdp-org/gd"
	"github.com/gdp-org/gd/net/dhttp"
	"github.com/gdp-org/gd/runtime/inject"
	"github.com/gin-gonic/gin"
	"hydra/app/api"
	"sync"
)

var (
	initOnce sync.Once
)

func Register(e *gd.Engine) {
	inject.RegisterOrFail("httpServerInit",func(g *gin.Engine) error {
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
