package main

import (
	"encoding/json"
	"strconv"
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
	key := "freightsrv-cep-region-" + cep
	// Save for one wekeend.
	_ = redisSet(key, region, time.Hour*168)
}

// Get CEP region.
func getCEPRegion(cep string) string {
	cep = strings.ReplaceAll(cep, "-", "")
	key := "freightsrv-cep-region-" + cep
	return redisGet(key)
}

//****************************************************************************
//	VIA CEP ADDRESS
//****************************************************************************
// Set via cep address.
func setViaCEPAddressCache(pCep *string, pAddressJson *string) {
	key := "freightsrv-via-cep-address-" + strings.ReplaceAll(*pCep, "-", "")
	_ = redisSet(key, string(*pAddressJson), time.Hour*2)
}

// Get via cep address.
func getViaCEPAddressCache(pCep *string) (*string, bool) {
	key := "freightsrv-via-cep-address-" + strings.ReplaceAll(*pCep, "-", "")
	addressJson := redisGet(key)
	if addressJson == "" {
		return nil, false
	}
	return &addressJson, true
}

//****************************************************************************
//	CORREIOS FREIGHTS
//****************************************************************************
func makeCorreiosKey(p *pack) string {
	return "freightsrv-correios-estimate-freight-" + strings.ReplaceAll(p.CEPOrigin, "-", "") + "-" + strings.ReplaceAll(p.CEPDestiny, "-", "") + "-" + strconv.Itoa(p.Weight) + "-" + strconv.Itoa(p.Length) + "-" + strconv.Itoa(p.Height) + "-" + strconv.Itoa(p.Width)
}

// Set Correios estimate delivery.
func setCorreiosCache(p *pack, frS []*freight) {
	frSJson, err := json.Marshal(frS)
	if checkError(err) {
		return
	}
	_ = redisSet(makeCorreiosKey(p), string(frSJson), time.Hour*2)
}

// Get Correios estimate delivery.
func getCorreiosCache(p *pack) (frS []*freight, ok bool) {
	frSJson := redisGet(makeCorreiosKey(p))
	// No key.
	if frSJson == "" {
		return frS, false
	}
	// log.Printf("frSJson: %s\n", frSJson)
	err := json.Unmarshal([]byte(frSJson), &frS)
	if checkError(err) {
		return frS, false
	}
	return frS, true
}
