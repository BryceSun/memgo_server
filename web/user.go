package web

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	router := Httprouter()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s! %v\n", ps.ByName("name"), r)
}
