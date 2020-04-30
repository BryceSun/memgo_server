package database

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v7"
)

const (
	usersKey  = "memgo:user:info:"
	loginsKey = "memgo:user:login"
	uidGenKey = "memgo:user:idg"
	userIndex = "memgo:user:index"
)

type UserInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Id       int64  `json:"id"`
}

func getZMembers(info *UserInfo) []*redis.Z {
	nameZ := &redis.Z{Score: float64(info.Id), Member: info.Name}
	mobileZ := &redis.Z{Score: float64(info.Id), Member: info.Mobile}
	emailZ := &redis.Z{Score: float64(info.Id), Member: info.Email}
	ms := []*redis.Z{nameZ, mobileZ, emailZ}
	return ms
}

//注册
func Register(info *UserInfo) (int64, error) {
	b, err := CanRegister(*info)
	if !b {
		return 0, err
	}

	uid, err := RedisClient.Incr(uidGenKey).Result()
	if err != nil {
		return 0, err
	}
	info.Id = uid
	infoJson, err := json.Marshal(*info)
	if err != nil {
		return 0, err
	}
	_, err = RedisClient.HSet(usersKey, uid, string(infoJson)).Result()
	if err != nil {
		return 0, err
	}
	_, err = RedisClient.ZAdd(userIndex, getZMembers(info)...).Result()
	if err != nil {
		return 0, err
	}
	return uid, nil
}

//查询已注册用户
func CanRegister(info UserInfo) (bool, error) {
	f, e := RedisClient.ZScore(userIndex, info.Name).Result()
	if e == nil && f > 0 {
		return false, errors.New("该用户名已被注册")
	}
	f, e = RedisClient.ZScore(userIndex, info.Email).Result()
	if e == nil && f > 0 {
		return false, errors.New("该邮箱已被注册")
	}
	f, e = RedisClient.ZScore(userIndex, info.Mobile).Result()
	if e == nil && f > 0 {
		return false, errors.New("该手机已被注册")
	}
	return true, nil
}

func login() {

}
