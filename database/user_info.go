package database

import (
	"encoding/json"
)

const (
	usersKey  = "memgo:users:"
	loginsKey = "memgo:logins"
	uidGenKey = "memgo:idgen"
)

type UserInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	id       int64  `json:"id"`
}

//注册
func Register(info *UserInfo) (int64, error) {
	uid, err := RedisClient.Incr(uidGenKey).Result()
	if err != nil {
		return 0, err
	}
	info.id = uid
	infoJson, err := json.Marshal(*info)
	if err != nil {
		return 0, err
	}
	_, err = RedisClient.HSet(usersKey, string(uid), string(infoJson)).Result()
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func login() {

}
