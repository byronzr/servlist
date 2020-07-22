package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/byronzr/servlist"
)

func main() {
	ip, err := servlist.Get("project_name_undefined")
	if err != nil {
		panic(err)
	}
	fmt.Println("get ip for: ", ip)
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
