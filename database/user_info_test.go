package database

import (
	"testing"
)

type userTest struct {
	UserInfo
	expect bool
}

var users = []userTest{
	{UserInfo: UserInfo{Name: "tony", Email: "tony@foxmail.com", Mobile: "15818376547", Password: "123456"}, expect: true},
	{UserInfo: UserInfo{Name: "jenny", Email: "jenny@foxmail.com", Mobile: "15818386548", Password: "123456"}, expect: true},
}

func TestRegister(t *testing.T) {

}

func TestCanRegister(t *testing.T) {
	user := UserInfo{Name: "wjw", Email: "qingyunxi@foxmail.com", Mobile: "15818366547", Password: "123456"}
	b, e := CanRegister(user)
	if b {
		t.Error("should' be false")
	}
	if e == nil {
		t.Error("error should be exist")
	}
}
