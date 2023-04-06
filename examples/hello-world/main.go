package main

import (
	"fmt"
	"github.com/lukasmwerner/pine"
	"github.com/lukasmwerner/pine/middlewares"
	"log"
	"net/http"
	"time"
)

func main() {
	r := pine.New()
	r.Use(middlewares.HTTPLogger())
	r.Use(middlewares.STolinskiTiming(time.Second, time.Second/2))
	r.HandleFunc("/hello/home", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprint("welcome home!")))
	})
	r.HandleFunc("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("Hello, %s", pine.Var(r, "name"))))
	})
	log.Fatalln(http.ListenAndServe(":8080", r))
}
