/**
 * Copyright 2020 gd-demo Author. All rights reserved.
 * Author: Chuck1024
 */

package boot

import (
	"github.com/chuck1024/gd"
	"github.com/chuck1024/gd/databases/redisdb"
	"github.com/chuck1024/gd/dlog"
	"github.com/chuck1024/gd/runtime/inject"
	"hydra/app/model"
	"hydra/app/service/sp"
)

func Inject(d *gd.Engine) {
	// inject SessionCache
	inject.Reg("UidCache", (*model.UidCache)(&model.UidCache{RedisConfig: &redisdb.RedisConfig{
		Addrs: d.Config("Redis", "addr").Strings(","),
	}}))

	// inject dependency
	inject.RegisterOrFail("serviceProvider", (*sp.ServiceProvider)(nil))
	err := sp.Init()
	if err != nil {
		dlog.Crashf("init package sp fail,err=%v", err)
	}
}
