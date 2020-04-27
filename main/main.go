package main

import (
	"github.com/memgo_server/web"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", web.Httprouter()))
}
