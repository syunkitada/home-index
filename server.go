package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	fmt.Println("Listen 0.0.0.0:3000")
	http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if strings.HasPrefix(r.URL.Path, "/") {
			http.ServeFile(w, r, r.URL.Path[1:])
		} else {
			http.NotFound(w, r)
		}
	}))
}
