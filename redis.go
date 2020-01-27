package main

import (
	"github.com/go-redis/redis/v7"
)

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

func redisSet(key string, val string) error {
	err := redisClient.Set(key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
