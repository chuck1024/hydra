/**
 * Copyright 2018 Author. All rights reserved.
 * Author: Chuck1024
 */

package cache

import (
	"errors"
	"fmt"
	"github.com/chuck1024/godog"
	"github.com/chuck1024/godog/store/cache"
	"github.com/chuck1024/godog/utils"
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

func getUuidKey(uuid uint64) string {
	return fmt.Sprintf("%s:%d", uuidPrefix, uuid)
}

func getPushKey(seq string) string {
	return fmt.Sprintf("%s:%s", pushPrefix, seq)
}

func SetUuid(uuid uint64) error {
	key := getUuidKey(uuid)
	value := utils.GetLocalIP()
	godog.Debug("[SetUuid] key: %s value:%s", key, value)

	err := cache.SetEx(key, expireTime, value)
	if err != nil {
		godog.Error("[SetUuid] redis SetEx occur error: %s, key:%s", err, key)
		return err
	}

	godog.Debug("[SetUuid] set key success. key:%s ", key)

	return nil
}

func GetUuid(uuid uint64) (value string, err error) {
	key := getUuidKey(uuid)
	godog.Debug("[GetUuid] key:%s ", key)

	value, err = cache.Get(key)
	if err != nil {
		newErr := fmt.Sprintf("%s", err)
		if newErr == "nil reply" {
			godog.Debug("[GetUuid] get value keyNotExist. key: %s", key)
			return "", KeyNotExist
		}

		godog.Error("[GetUuid] redis Get occur error: %s, key:%s", err, key)
		return
	}

	if len(value) == 0 {
		return "", KeyNotExist
	}

	godog.Debug("[GetUuid] localAddr: %s", value)

	return
}

func DelUuid(uuid uint64) error {
	key := getUuidKey(uuid)
	godog.Debug("[DelUuid] key: %s", key)

	_, err := cache.Del(key)
	if err != nil {
		godog.Error("[DelUuid] redis Del occur error: %s, key:%s", err, key)
		return err
	}

	godog.Debug("[DelUuid] Del key success. key:%s ", key)

	return nil
}

func SetPush(id string) error {
	key := getPushKey(id)
	godog.Debug("[SetPush] key: %s", key)

	err := cache.SetEx(key, pushExpireTime, "check")
	if err != nil {
		godog.Error("[SetPush] redis setEx occur error: %s, key:%s", err, key)
		return err
	}

	godog.Debug("[SetPush] set key success. key: %s", key)

	return nil
}

func GetPush(id string) bool {
	key := getPushKey(id)
	godog.Debug("[GetPush] key:%s ", key)

	value, err := cache.Get(key)
	if err != nil {
		newErr := fmt.Sprintf("%s", err)
		if newErr == "nil reply" {
			godog.Debug("[GetPush] redis check no repeat key: %s ", key)
			return false
		}

		godog.Error("[GetPush] redis Get occur error : %s, key:%s", err, key)
		return false
	}

	if len(value) == 0 {
		return false
	}

	return true
}
