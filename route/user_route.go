package route

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	. "github.com/memgo_server/databean"
	. "github.com/memgo_server/handler"
	"net/http"
	"strconv"
)

func init() {
	router := HttpRouter()
	router.PUT("/user/account", register)
	router.PUT("/user/token", login)
	router.DELETE("/user/token", logout)
}

func register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	mobile := r.FormValue("mobile")
	user := &UserInfo{
		Username: username,
		Password: password,
		Email:    email,
		Mobile:   mobile,
	}
	_, e := Register(user)
	if e != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		i := ErrorInfo{http.StatusNotAcceptable, "昵称密码不可为空"}
		encoder := json.NewEncoder(w)
		encoder.Encode(i)
	}
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := UserInfo{Username: username, Password: password}
	token, e := Login(&user)
	if e != nil {
		// todo
		w.Write([]byte(e.Error()))
	}
	w.Header()["authentication"] = []string{"Bearer" + token}
}

func logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := r.FormValue("id")
	if len(id) != 0 {
		uid, e := strconv.ParseInt(id, 10, 64)
		if e != nil {
			// todo
			w.Write([]byte(e.Error()))
		}
		e = Logout(uid)
		if e != nil {
			// todo
			w.Write([]byte(e.Error()))
		}
	}
}
