/**
 * Copyright 2020 gd-demo Author. All rights reserved.
 * Author: Chuck1024
 */

package boot

import (
	"github.com/gdp-org/gd"
	"github.com/gdp-org/gd/databases/redisdb"
	"github.com/gdp-org/gd/runtime/inject"
	"hydra/app/service/sp"
)

func Inject() {
	// inject redisClient and init redis pool client
	inject.RegisterOrFail("redisClient", (*redisdb.RedisPoolClient)(&redisdb.RedisPoolClient{
		PoolName: "hydra",
	}))

	// inject dependency
	inject.RegisterOrFail("serviceProvider", (*sp.ServiceProvider)(nil))
	err := sp.Init()
	if err != nil {
		gd.Crashf("init package sp fail,err=%v", err)
	}
}
