package database

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	. "github.com/memgo_server/databean"
	"strconv"
)

const (
	userBaseKey  = memgoKey + ":user"
	userInfosKey = userBaseKey + ":info"
	logExpireKey = userBaseKey + ":login"
	userIdGenKey = userBaseKey + ":idgen"
	userIdIndex  = userBaseKey + ":index"
)

func ClearUser() (int64, error) {
	pattern := userBaseKey + "*"
	return clearRedis(pattern)
}

func getZMembers(info *UserInfo) []*redis.Z {
	nameZ := &redis.Z{Score: float64(info.Id), Member: info.Username}
	mobileZ := &redis.Z{Score: float64(info.Id), Member: info.Mobile}
	emailZ := &redis.Z{Score: float64(info.Id), Member: info.Email}
	ms := []*redis.Z{nameZ, mobileZ, emailZ}
	return ms
}

//注册
func AddUser(user *UserInfo) (int64, error) {
	uid, err := RedisClient.Incr(userIdGenKey).Result()
	if err != nil {
		return 0, err
	}
	user.Id = uid

	infoJson, err := json.Marshal(*user)
	if err != nil {
		return 0, err
	}
	_, err = RedisClient.HSet(userInfosKey, uid, string(infoJson)).Result()
	if err != nil {
		return 0, err
	}
	_, err = RedisClient.ZAdd(userIdIndex, getZMembers(user)...).Result()
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func GetUserId(filedValue string) (int64, error) {
	uid, e := RedisClient.ZScore(userIdIndex, filedValue).Result()
	//if e != nil{
	//	return 0,e
	//}
	return int64(uid), e
}

func GetUserInfo(uid int64) (user UserInfo, err error) {
	// get user info by id
	userJsonStr, err := RedisClient.HGet(userInfosKey, fmt.Sprint(uid)).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(userJsonStr), &user)
	return
}

func LogExpireTime(uid, expireTime int64) error {
	_, err := RedisClient.HSet(logExpireKey, uid, expireTime).Result()
	return err
}

func GetExpireTime(uid int64) (int64, error) {
	time, err := RedisClient.HGet(logExpireKey, fmt.Sprint(uid)).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(time, 10, 64)
}
