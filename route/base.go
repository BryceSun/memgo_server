package route

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

var (
	router *httprouter.Router
)

type HH httprouter.Handle

func (h HH) timeCost() HH {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		start := time.Now()
		h(w, r, p)
		log.Println("It cost v% to login", time.Since(start))
	}
}

func (h HH) auth() HH {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Println("auth before")
		h(w, r, p)
		log.Println("auth after")
	}
}

func (h HH) withFilter() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		h.timeCost().auth()
	}
}

func init() {
	router = httprouter.New()
}

type ErrorInfo struct {
	Code int    `json:"code"` //字段首字母要大写才能导出
	Msg  string `json:"msg"`
}

func HttpRouter() *httprouter.Router {
	return router
}
