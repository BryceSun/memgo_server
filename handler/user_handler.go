package handler

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/memgo_server/database"
	. "github.com/memgo_server/databean"
	"time"
)

var (
	UsernameHaveRegistered = errors.New("该用户名已被注册")
	EmailHaveRegistered    = errors.New("该邮箱已被注册")
	MobileHaveRegistered   = errors.New("该手机号已被注册")
	CanNotBeAllEmpty       = errors.New("登录名不可为空")
	PasswordIsRequired     = errors.New("密码不可为空")
	PasswordIsWrong        = errors.New("密码错误")
	ParamsAreEmpty         = errors.New("缺乏有效参数")
	UserInfoNotExist       = errors.New("用户不存在")
	NotLogonUser           = errors.New("非在线用户")
)

//查询是否可以注册
func canRegister(user *UserInfo) (bool, error) {
	uid, e := database.GetUserId(user.Username)
	if e == nil && uid > 0 {
		return false, UsernameHaveRegistered
	}
	uid, e = database.GetUserId(user.Email)
	if e == nil && uid > 0 {
		return false, EmailHaveRegistered
	}
	uid, e = database.GetUserId(user.Mobile)
	if e == nil && uid > 0 {
		return false, MobileHaveRegistered
	}
	return true, nil
}

// 注册
func Register(user *UserInfo) (int64, error) {
	if len(user.Password) == 0 || len(user.Username) == 0 {
		return 0, ParamsAreEmpty
	}
	allow, e := canRegister(user)
	if !allow {
		return 0, e
	}
	user.Password = fmt.Sprintf("%x", md5.Sum([]byte(user.Password)))
	return database.AddUser(user)
}

// 登录
func Login(user Loginer) (auth string, err error) {
	loginName := user.GetName()
	if 0 == len(loginName) {
		err = CanNotBeAllEmpty
		return
	}
	password := user.GetPassword()
	if 0 == len(password) {
		err = PasswordIsRequired
		return
	}
	// get user id by login name
	uid, err := database.GetUserId(loginName)
	if err != nil {
		err = UserInfoNotExist
		return
	}
	// get user info by id
	u, err := database.GetUserInfo(uid)
	if err != nil {
		return
	}
	// check password
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	if password != u.Password {
		err = PasswordIsWrong
		return
	}
	return getToken(&u)
}

func getToken(user *UserInfo) (string, error) {

	mySigningKey := []byte("MemGoIsBestToolToUse")

	expireTime := time.Now().Add(24 * time.Second).Unix()
	e := database.LogExpireTime(user.Id, expireTime)
	if e != nil {
		return "", e
	}

	claims := make(jwt.MapClaims)
	claims["exp"] = expireTime
	claims["iat"] = time.Now().Unix()
	claims["uid"] = user.Id
	claims["unm"] = user.Username
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func Logout(uid int64) error {
	now := time.Now().Unix()
	exTime, e := database.GetExpireTime(uid)
	if e != nil {
		return NotLogonUser
	}
	if exTime > now {
		return database.LogExpireTime(uid, now)
	}
	return nil
}
