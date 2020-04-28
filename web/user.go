package web

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
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
	if len(username) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		e := errorInfo{http.StatusNotAcceptable, "昵称密码不可为空"}
		fmt.Fprint(w, e)
	} else {
		fmt.Fprint(w, "Welcome!\n")
	}
}

func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s! %v\n", ps.ByName("name"), r)
}

func logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s! %v\n", ps.ByName("name"), r)
}
