package main

import (
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
)

// Get.
func redisGet(key string) string {
	val, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		return ""
		// log.Printf("Key not exist")
	} else if err != nil {
		checkError(err)
		return ""
	} else {
		return val
		// log.Printf("Val: %+v\n", val)
	}
}

// Set.
func redisSet(key string, val string, exp time.Duration) error {
	err := redisClient.Set(key, val, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

// Del.
func redisDel(key string) error {
	return redisClient.Del(key).Err()
}

//****************************************************************************
//	CEP
//****************************************************************************
// Set CEP region.
func setCEPRegion(cep string, region string) {
	cep = strings.ReplaceAll(cep, "-", "")
	key := "cep-region-" + cep
	// Save for one wekeend.
	_ = redisSet(key, region, time.Hour*168)
}

// Get CEP region.
func getCEPRegion(cep string) string {
	cep = strings.ReplaceAll(cep, "-", "")
	key := "cep-region-" + cep
	return redisGet(key)
}
