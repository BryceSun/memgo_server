package database

import (
	. "github.com/memgo_server/databean"
	"reflect"
	"testing"
)

var users = []UserInfo{
	{
		Username: "Tom",
		Email:    "15818366548@gmail.com",
		Mobile:   "15717366548",
		Password: "123456",
	},
	{
		Username: "Jack",
		Email:    "15818366549@gmail.com",
		Mobile:   "15717366549",
		Password: "123456",
	},
	{
		Username: "Jenny",
		Email:    "15818366578@gmail.com",
		Mobile:   "15717366578",
		Password: "123456",
	},
}

var clear = true

func TestAddUser(t *testing.T) {
	_, e := clearUser()
	if e != nil {
		t.Errorf("clearing user redis has error:%v ", e)
	}

	for i, user := range users {
		uid, e := AddUser(&user)
		if e != nil {
			t.Errorf("have error:%v", e)
		}
		if uid != int64(i+1) {
			t.Errorf("uid should be %v", uid)
		}
	}

	if clear {
		i, e := clearUser()
		if e != nil {
			t.Errorf("clearing user redis has error:%v ", e)
		}
		if int(i) != len(users) {
			t.Errorf("%v keys should be cleared", len(users))
		}
	}
}

func TestGetUserId(t *testing.T) {
	_, e := clearUser()
	if e != nil {
		t.Errorf("clearing user redis has error:%v ", e)
	}
	clear = false
	TestAddUser(t)
	for i, user := range users {
		uid, err := GetUserId(user.Username)
		if err != nil {
			t.Errorf("getting user id has error")
		}
		if int(uid) != i+1 {
			t.Errorf("uid should be %v", uid)
		}
		uid, err = GetUserId(user.Email)
		if err != nil {
			t.Errorf("getting user id has error")
		}
		if int(uid) != i+1 {
			t.Errorf("uid should be %v", uid)
		}
		uid, err = GetUserId(user.Mobile)
		if err != nil {
			t.Errorf("getting user id has error")
		}
		if int(uid) != i+1 {
			t.Errorf("uid should be %v", uid)
		}
	}
}

func TestGetUserInfo(t *testing.T) {
	clear = false
	TestAddUser(t)
	for i, user := range users {
		uid, _ := GetUserId(user.Username)
		user.Id = int64(i + 1)
		guser, _ := GetUserInfo(uid)
		if !reflect.DeepEqual(user, guser) {
			t.Error("guser should be same as user")
		}
	}
}
