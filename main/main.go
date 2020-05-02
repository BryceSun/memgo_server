package main

import (
	"github.com/memgo_server/route"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", route.HttpRouter()))
}
