/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package cache

import (
	"errors"
	"fmt"
	"github.com/chuck1024/doglog"
	"github.com/chuck1024/godog"
	"github.com/chuck1024/godog/utils"
	"github.com/chuck1024/redisdb"
)

// conf.json redis tag : redis
const redisTag = "redis"

const (
	uuidPrefix     = "hydra:uuid" // key:uuid value:localAddr
	expireTime     = 120
	pushPrefix     = "hydra:push" // key:seq
	pushExpireTime = 3600
)

var (
	redisHandle *redisdb.RedisPool
	KeyNotExist = errors.New("KeyNotExist")
)

func Init(dog *godog.Engine) {
	url, _ := dog.Config.String(redisTag)
	cfg, err := redisdb.RedisConfigFromURLString(url)
	if err != nil {
		doglog.Error("redis Init RedisConfigFromURLString occur error:%s", err)
		return
	}

	redisHandle, err = redisdb.NewRedisPools(cfg)
	if err != nil {
		doglog.Error("redis Init NewRedisPools occur error:%s", err)
		return
	}
}

func getUuidKey(uuid uint64) string {
	return fmt.Sprintf("%s:%d", uuidPrefix, uuid)
}

func getPushKey(seq string) string {
	return fmt.Sprintf("%s:%s", pushPrefix, seq)
}

func SetUuid(uuid uint64) error {
	key := getUuidKey(uuid)
	value := utils.GetLocalIP()
	//value := utils.GetLocalIP() + ":" + strconv.Itoa(doglog.AppConfig.BaseConfig.Server.HttpPort)
	doglog.Debug("[SetUuid] key: %s value:%s", key, value)

	err := redisHandle.SetEx(key, expireTime, value)
	if err != nil {
		doglog.Error("[SetUuid] redis SetEx occur error: %s, key:%s", err, key)
		return err
	}

	doglog.Debug("[SetUuid] set key success. key:%s ", key)

	return nil
}

func GetUuid(uuid uint64) (value string, err error) {
	key := getUuidKey(uuid)
	doglog.Debug("[GetUuid] key:%s ", key)

	value, err = redisHandle.Get(key)
	if err != nil {
		newErr := fmt.Sprintf("%s", err)
		if newErr == "nil reply" {
			doglog.Debug("[GetUuid] get value keyNotExist. key: %s", key)
			return "", KeyNotExist
		}

		doglog.Error("[GetUuid] redis Get occur error: %s, key:%s", err, key)
		return
	}

	if len(value) == 0 {
		return "", KeyNotExist
	}

	doglog.Debug("[GetUuid] localAddr: %s", value)

	return
}

func DelUuid(uuid uint64) error {
	key := getUuidKey(uuid)
	doglog.Debug("[DelUuid] key: %s", key)

	err := redisHandle.Del(key)
	if err != nil {
		doglog.Error("[DelUuid] redis Del occur error: %s, key:%s", err, key)
		return err
	}

	doglog.Debug("[DelUuid] Del key success. key:%s ", key)

	return nil
}

func SetPush(id string) error {
	key := getPushKey(id)
	doglog.Debug("[SetPush] key: %s", key)

	err := redisHandle.SetEx(key, pushExpireTime, "check")
	if err != nil {
		doglog.Error("[SetPush] redis setEx occur error: %s, key:%s", err, key)
		return err
	}

	doglog.Debug("[SetPush] set key success. key: %s", key)

	return nil
}

func GetPush(id string) bool {
	key := getPushKey(id)
	doglog.Debug("[GetPush] key:%s ", key)

	value, err := redisHandle.Get(key)
	if err != nil {
		newErr := fmt.Sprintf("%s", err)
		if newErr == "nil reply" {
			doglog.Debug("[GetPush] redis check no repeat key: %s ", key)
			return false
		}

		doglog.Error("[GetPush] redis Get occur error : %s, key:%s", err, key)
		return false
	}

	if len(value) == 0 {
		return false
	}

	return true
}
