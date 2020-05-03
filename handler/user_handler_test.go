package handler

import (
	"github.com/memgo_server/database"
	. "github.com/memgo_server/databean"
	"log"
	"testing"
)

type userTest struct {
	UserInfo
	IdExpect    int64
	ErrorExpect error
}

var user1 = UserInfo{Username: "tony", Email: "tony@foxmail.com", Mobile: "15818376547", Password: "123456"}
var user2 = UserInfo{Username: "tony", Email: "jenny@foxmail.com", Mobile: "15818386548", Password: "123456"}
var user3 = UserInfo{Username: "tom", Email: "tony@foxmail.com", Mobile: "15818386549", Password: "123456"}
var user4 = UserInfo{Username: "jack", Email: "jack@foxmail.com", Mobile: "15818376547", Password: "123456"}
var user5 = UserInfo{Username: "bryce", Email: "bryce@foxmail.com", Mobile: "15016724757", Password: "123456"}

var user1login1 = UserInfo{}
var user1login2 = UserInfo{Password: "123456"}
var user1login3 = UserInfo{Username: "tony"}
var user1login4 = UserInfo{Email: "tony@foxmail.com"}
var user1login5 = UserInfo{Mobile: "15818376547"}
var user1login6 = UserInfo{Username: "tony", Password: "123455"}
var user1login7 = UserInfo{Email: "tony@foxmail.com", Password: "123455"}
var user1login8 = UserInfo{Mobile: "15818376547", Password: "123455"}
var user1login14 = UserInfo{Username: "tom", Password: "123455"}
var user1login15 = UserInfo{Email: "tom@foxmail.com", Password: "123455"}
var user1login16 = UserInfo{Mobile: "15818376548", Password: "123455"}
var user1login9 = UserInfo{Username: "tony", Password: "123456"}
var user1login10 = UserInfo{Email: "tony@foxmail.com", Password: "123456"}
var user1login11 = UserInfo{Mobile: "15818376547", Password: "123456"}
var user1login12 = UserInfo{Username: "tony", Email: "tony@foxmail.com", Mobile: "15818376547", Password: "123455"}
var user1login13 = UserInfo{Username: "tony", Email: "tony@foxmail.com", Mobile: "15818376547", Password: "123456"}
var user5login1 = UserInfo{Username: "bryce", Email: "bryce@foxmail.com", Mobile: "15016724757", Password: "123456"}

var userTests = []userTest{
	{UserInfo: user1, IdExpect: 1, ErrorExpect: nil},
	{UserInfo: user2, IdExpect: 0, ErrorExpect: UsernameHaveRegistered},
	{UserInfo: user3, IdExpect: 0, ErrorExpect: EmailHaveRegistered},
	{UserInfo: user4, IdExpect: 0, ErrorExpect: MobileHaveRegistered},
	{UserInfo: user5, IdExpect: 2, ErrorExpect: nil},
}

var user1LoginTests = []userTest{
	{UserInfo: user1login1, ErrorExpect: CanNotBeAllEmpty},
	{UserInfo: user1login2, ErrorExpect: CanNotBeAllEmpty},
	{UserInfo: user1login3, ErrorExpect: PasswordIsRequired},
	{UserInfo: user1login4, ErrorExpect: PasswordIsRequired},
	{UserInfo: user1login5, ErrorExpect: PasswordIsRequired},
	{UserInfo: user1login6, ErrorExpect: PasswordIsWrong},
	{UserInfo: user1login7, ErrorExpect: PasswordIsWrong},
	{UserInfo: user1login8, ErrorExpect: PasswordIsWrong},
	{UserInfo: user1login14, ErrorExpect: UserInfoNotExist},
	{UserInfo: user1login15, ErrorExpect: UserInfoNotExist},
	{UserInfo: user1login16, ErrorExpect: UserInfoNotExist},
	{UserInfo: user1login9, ErrorExpect: nil},
	{UserInfo: user1login10, ErrorExpect: nil},
	{UserInfo: user1login11, ErrorExpect: nil},
	{UserInfo: user1login12, ErrorExpect: PasswordIsWrong},
	{UserInfo: user1login13, ErrorExpect: nil},
	{UserInfo: user5login1, ErrorExpect: nil},
}

func TestRegister(t *testing.T) {
	database.ClearUser()
	for _, u := range userTests {
		ri, err := Register(&u.UserInfo)
		if ri != u.IdExpect {
			t.Errorf("id should be %v", u.IdExpect)
		}
		if err != u.ErrorExpect {
			t.Errorf("error shoul be %v", u.ErrorExpect.Error())
		}
	}
}

func TestLogin(t *testing.T) {
	database.ClearUser()
	Register(&user1)
	Register(&user5)
	for _, u := range user1LoginTests {
		//time.Sleep(time.Second/2)
		token, e := Login(&u.UserInfo)
		if e != nil && e != u.ErrorExpect {
			t.Errorf("error happened,but not expected :%v", e)
			continue
		}
		if e == nil && len(token) == 0 {
			t.Error("token must not be empty")
		}
		if len(token) != 0 {
			log.Println(token)
		}
	}
}

func TestLogout(t *testing.T) {
	TestLogin(t)
	e := Logout(1)
	if e != nil {
		t.Errorf("logout error:%v", e)
	}
	e = Logout(2)
	if e != nil {
		t.Errorf("logout error:%v", e)
	}
	e = Logout(3)
	if e != NotLogonUser {
		t.Error("should be expected error:", e)
	}
}
