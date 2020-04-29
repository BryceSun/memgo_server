package web

import "github.com/julienschmidt/httprouter"

var router *httprouter.Router

func init() {
	router = httprouter.New()
}

type ErrorInfo struct {
	Code int    `json:"code"` //字段首字母要大写才能导出
	Msg  string `json:"msg"`
}

func Httprouter() *httprouter.Router {
	return router
}
