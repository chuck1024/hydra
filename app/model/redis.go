/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package model

import (
	"errors"
	"fmt"
	"github.com/chuck1024/gd/databases/redisdb"
	"github.com/chuck1024/gd/dlog"
	"github.com/chuck1024/gd/utls/network"
)

const (
	uuidPrefix     = "hydra:uuid" // key:uuid value:localAddr
	expireTime     = 120
	pushPrefix     = "hydra:push" // key:seq
	pushExpireTime = 3600
)

var (
	KeyNotExist = errors.New("KeyNotExist")
)

type UidCache struct {
	RedisConfig *redisdb.RedisConfig
	redisPool   *redisdb.RedisPool
}

func (c *UidCache) Start() (err error) {
	c.redisPool, err = redisdb.NewRedisPools(c.RedisConfig)
	if err != nil {
		return
	}
	return nil
}

func getUuidKey(uuid uint64) string {
	return fmt.Sprintf("%s:%d", uuidPrefix, uuid)
}

func getPushKey(seq string) string {
	return fmt.Sprintf("%s:%s", pushPrefix, seq)
}

func (c *UidCache) SetUuid(uuid uint64) error {
	key := getUuidKey(uuid)
	value := network.GetLocalIP()
	//value := utils.GetLocalIP() + ":" + strconv.Itoa(dlog.AppConfig.BaseConfig.Server.HttpPort)
	dlog.Debug("[SetUuid] key: %s value:%s", key, value)

	err := c.redisPool.SetEx(key, expireTime, value)
	if err != nil {
		dlog.Error("[SetUuid] redis SetEx occur error: %s, key:%s", err, key)
		return err
	}

	dlog.Debug("[SetUuid] set key success. key:%s ", key)

	return nil
}

func (c *UidCache) GetUuid(uuid uint64) (value string, err error) {
	key := getUuidKey(uuid)
	dlog.Debug("[GetUuid] key:%s ", key)

	value, err = c.redisPool.Get(key)
	if err != nil {
		newErr := fmt.Sprintf("%s", err)
		if newErr == "nil reply" {
			dlog.Debug("[GetUuid] get value keyNotExist. key: %s", key)
			return "", KeyNotExist
		}

		dlog.Error("[GetUuid] redis Get occur error: %s, key:%s", err, key)
		return
	}

	if len(value) == 0 {
		return "", KeyNotExist
	}

	dlog.Debug("[GetUuid] localAddr: %s", value)

	return
}

func (c *UidCache) DelUuid(uuid uint64) error {
	key := getUuidKey(uuid)
	dlog.Debug("[DelUuid] key: %s", key)

	err := c.redisPool.Del(key)
	if err != nil {
		dlog.Error("[DelUuid] redis Del occur error: %s, key:%s", err, key)
		return err
	}

	dlog.Debug("[DelUuid] Del key success. key:%s ", key)

	return nil
}

func (c *UidCache) SetPush(id string) error {
	key := getPushKey(id)
	dlog.Debug("[SetPush] key: %s", key)

	err := c.redisPool.SetEx(key, pushExpireTime, "check")
	if err != nil {
		dlog.Error("[SetPush] redis setEx occur error: %s, key:%s", err, key)
		return err
	}

	dlog.Debug("[SetPush] set key success. key: %s", key)

	return nil
}

func (c *UidCache) GetPush(id string) bool {
	key := getPushKey(id)
	dlog.Debug("[GetPush] key:%s ", key)

	value, err := c.redisPool.Get(key)
	if err != nil {
		newErr := fmt.Sprintf("%s", err)
		if newErr == "nil reply" {
			dlog.Debug("[GetPush] redis check no repeat key: %s ", key)
			return false
		}

		dlog.Error("[GetPush] redis Get occur error : %s, key:%s", err, key)
		return false
	}

	if len(value) == 0 {
		return false
	}

	return true
}
