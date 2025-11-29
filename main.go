package main

import (
	"fmt"
	"net/http"
	"sync"

	"imga/api"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func main() {
	urlMap := make(map[string]string)
	var urlMapMu sync.RWMutex

	http.HandleFunc("/shorten", api.CreateShortURLHandler(urlMap, &urlMapMu))
	http.HandleFunc("/{shortCode}", api.RedirectToOriginalURLHandler(urlMap, &urlMapMu))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
