package main

import (
	"fmt"
	"github.com/lukasmwerner/pine"
	"log"
	"net/http"
)

func main() {
	r := pine.New()
	r.Handle("/hello/home", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprint("welcome home!")))
	})
	r.Handle("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("Hello, %s", pine.Var(r, "name"))))
	})
	log.Fatalln(http.ListenAndServe(":8080", r))
}
