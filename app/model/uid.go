/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package model

import (
	"errors"
	"fmt"
	"github.com/gdp-org/gd"
	"github.com/gdp-org/gd/databases/redisdb"
	"github.com/gdp-org/gd/utls/network"
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
	RedisClient *redisdb.RedisPoolClient `inject:"redisClient"`
}

func (c *UidCache) Start() (err error) {
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
	gd.Debug("[SetUuid] key: %s value:%s", key, value)

	err := c.RedisClient.SetEx(key, expireTime, value)
	if err != nil {
		gd.Error("[SetUuid] redis SetEx occur error: %s, key:%s", err, key)
		return err
	}

	gd.Debug("[SetUuid] set key success. key:%s ", key)

	return nil
}

func (c *UidCache) GetUuid(uuid uint64) (value string, err error) {
	key := getUuidKey(uuid)
	gd.Debug("[GetUuid] key:%s ", key)

	value, err = c.RedisClient.Get(key)
	if err != nil {
		newErr := fmt.Sprintf("%s", err)
		if newErr == "nil reply" {
			gd.Debug("[GetUuid] get value keyNotExist. key: %s", key)
			return "", KeyNotExist
		}

		gd.Error("[GetUuid] redis Get occur error: %s, key:%s", err, key)
		return
	}

	if len(value) == 0 {
		return "", KeyNotExist
	}

	gd.Debug("[GetUuid] localAddr: %s", value)

	return
}

func (c *UidCache) DelUuid(uuid uint64) error {
	key := getUuidKey(uuid)
	gd.Debug("[DelUuid] key: %s", key)

	err := c.RedisClient.Del(key)
	if err != nil {
		gd.Error("[DelUuid] redis Del occur error: %s, key:%s", err, key)
		return err
	}

	gd.Debug("[DelUuid] Del key success. key:%s ", key)

	return nil
}

func (c *UidCache) SetPush(id string) error {
	key := getPushKey(id)
	gd.Debug("[SetPush] key: %s", key)

	err := c.RedisClient.SetEx(key, pushExpireTime, "check")
	if err != nil {
		gd.Error("[SetPush] redis setEx occur error: %s, key:%s", err, key)
		return err
	}

	gd.Debug("[SetPush] set key success. key: %s", key)

	return nil
}

func (c *UidCache) GetPush(id string) bool {
	key := getPushKey(id)
	gd.Debug("[GetPush] key:%s ", key)

	value, err := c.RedisClient.Get(key)
	if err != nil {
		newErr := fmt.Sprintf("%s", err)
		if newErr == "nil reply" {
			gd.Debug("[GetPush] redis check no repeat key: %s ", key)
			return false
		}

		gd.Error("[GetPush] redis Get occur error : %s, key:%s", err, key)
		return false
	}

	if len(value) == 0 {
		return false
	}

	return true
}
