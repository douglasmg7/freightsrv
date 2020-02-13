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

//****************************************************************************
//	VIA CEP ADDRESS
//****************************************************************************
// Set via cep address.
func setViaCEPAddressCache(pCep *string, pAddressJson *string) {
	key := "via-cep-address-" + strings.ReplaceAll(*pCep, "-", "")
	_ = redisSet(key, string(*pAddressJson), time.Hour*2)
}

// Get via cep address.
func getViaCEPAddressCache(pCep *string) (*string, bool) {
	key := "via-cep-address-" + strings.ReplaceAll(*pCep, "-", "")
	addressJson := redisGet(key)
	if addressJson == "" {
		return nil, false
	}
	return &addressJson, true
}

// // Set via cep address.
// func setViaCEPAddressCache(cep string, address viaCEPAddress) {
// cep = strings.ReplaceAll(cep, "-", "")
// key := "via-cep-address-" + cep
// addressJson, err := json.Marshal(address)
// if checkError(err) {
// return
// }
// _ = redisSet(key, string(addressJson), time.Hour*2)
// }

// // Get via cep address.
// func getViaCEPAddressCache(cep string) (pAddress *ViaCEPAddress, ok bool) {
// pAddress = &viaCEPAddress{}
// cep = strings.ReplaceAll(cep, "-", "")
// key := "via-cep-address-" + cep
// addressJson := redisGet(key)
// err := json.Unmarshal([]byte(addressJson), &pAddress)
// if checkError(err) {
// return pAddress, false
// }
// return pAddress, true
// }
