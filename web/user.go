package web

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/memgo_server/database"
	"net/http"
)

func init() {
	router := Httprouter()
	router.PUT("/user/account", register)
	router.PUT("/user/token", login)
	router.DELETE("/user/token", logout)
}

func register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	if len(username) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		i := ErrorInfo{http.StatusNotAcceptable, "昵称密码不可为空"}
		encoder := json.NewEncoder(w)
		encoder.Encode(i)
	} else {
		user := &database.UserInfo{
			Name:     username,
			Password: password,
			Email:    email,
		}
		uid, err := database.Register(user)
		if err != nil {
			fmt.Fprint(w, "系统错误!\n")
		}
		fmt.Fprint(w, uid)
	}
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		i := ErrorInfo{http.StatusNotAcceptable, "昵称密码不可为空"}
		encoder := json.NewEncoder(w)
		encoder.Encode(i)
	} else {
		auth, err := database.RedisClient.Get(username).Result()
		if err != nil {
			fmt.Fprintf(w, "System error! %v \n", err)
			return
		}
		if auth == password {
			fmt.Fprint(w, "Welcome login memgo!\n")
		} else {
			fmt.Fprint(w, "password or username is wrong!\n")
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.FormValue("username")
	if len(username) == 0 {
		fmt.Fprint(w, "username is required!\n")
	}
	t, err := database.RedisClient.Del(username).Result()
	if err != nil {
		fmt.Fprintf(w, "System error! %v \n", err)
		return
	}
	if t == 1 {
		fmt.Fprint(w, "logout")
	}
}
