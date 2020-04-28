package web

import "github.com/julienschmidt/httprouter"

var router *httprouter.Router

func init() {
	router = httprouter.New()
}

type errorInfo struct {
	code int
	msg  string
}

func Httprouter() *httprouter.Router {
	return router
}
