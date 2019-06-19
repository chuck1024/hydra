/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package main

import (
	"github.com/chuck1024/doglog"
	"github.com/chuck1024/godog"
	"github.com/chuck1024/godog/net/httplib"
	"github.com/chuck1024/hydra/controller"
	"github.com/chuck1024/hydra/dao/cache"
	"github.com/chuck1024/hydra/service"
	"github.com/chuck1024/redisdb"
	"github.com/gin-gonic/gin"
)

func register(dog *godog.Engine) {
	dog.HttpServer.DefaultAddHandler("/push", controller.PushControl)
	dog.HttpServer.DefaultAddHandler("/query", controller.QueryControl)
	dog.HttpServer.DefaultAddHandler("/route", controller.RouteControl)

	dog.HttpServer.SetInit(func(g *gin.Engine) error {
		r := g.Group("")
		r.Use(
			httplib.GroupFilter(),
			httplib.Logger(),
		)

		for k, v := range dog.HttpServer.DefaultHandlerMap {
			f, err := httplib.Wrap(v)
			if err != nil {
				return err
			}
			r.GET(k, f)
			r.POST(k, f)
		}

		rr := g.Group("ws")
		rr.GET("/hydra", controller.WsControl)

		return nil
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	dog := godog.Default()

	url, _ := dog.Config.String("redis")
	cfg, err := redisdb.RedisConfigFromURLString(url)
	if err != nil{
		doglog.Error("redisdb.RedisConfigFromURLString occur error:%s", err)
		return
	}

	cache.RedisHandle, err = redisdb.NewRedisPools(cfg)
	if err != nil{
		doglog.Error("redisdb.NewRedisPools occur error:%s", err)
		return
	}

	register(dog)

	go service.Start()

	if err := dog.Run(); err != nil {
		doglog.Error("Error occurs, error = %s", err.Error())
		return
	}
}
