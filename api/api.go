package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Shorten URL params
type ShortenURLParams struct {
	URL string
}

// Shorten URL response
type ShortenURLResponse struct {
	ShortCode string
	ShortURL string
	OriginalURL string
}

type Error struct {
	Code int
	Message string
}

// URL validator
func isValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	return true
}


func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("invalid request type to /shorten")
		return
	}

	// Read and decode JSON body into ShortenUrlParams
	var params ShortenURLParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		if err == io.EOF {
			http.Error(w, "Empty request body", http.StatusBadRequest)
		} else {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		}
		log.Printf("error decoding /shorten request body: %v", err)
		return
	}

	// If URL is not valid, return error
	isValid := isValidURL(params.URL)
	if !isValid {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Printf("Invalid URL: %v", params.URL)
		return
	}

	log.Printf("shorten request for URL: %s", params.URL)

	fmt.Fprintf(w, "create short url route\n")
}