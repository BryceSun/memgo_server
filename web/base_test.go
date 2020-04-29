package web

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJson(t *testing.T) {
	type Foo struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}

	fm, _ := json.Marshal(Foo{Number: 1, Title: "test"})
	fmt.Print(string(fm)) // write response to ResponseWriter (w)
}
