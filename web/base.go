package web

import "github.com/julienschmidt/httprouter"

var router *httprouter.Router

func init() {
	router = httprouter.New()
}

func Httprouter() *httprouter.Router {
	return router
}
