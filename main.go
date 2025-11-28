package main

import (
	"fmt"
	"net/http"

	"imga/api"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func main() {
	http.HandleFunc("/shorten", api.CreateShortUrlHandler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}