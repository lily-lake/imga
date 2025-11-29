package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Shorten URL params
type ShortenUrlParams struct {
	Url string
}

// Shorten URL response
type ShortenUrlResponse struct {
	ShortCode string
	ShortUrl string
	OriginalUrl string
}

type Error struct {
	Code int
	Message string
}




func CreateShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("invalid request type to /shorten")
		return
	}

	// Read and decode JSON body into ShortenUrlParams
	var params ShortenUrlParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		if err == io.EOF {
			http.Error(w, "Empty request body", http.StatusBadRequest)
		} else {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		}
		log.Printf("error decoding /shorten request body: %v", err)
		return
	}

	log.Printf("shorten request for URL: %s", params.Url)

	fmt.Fprintf(w, "create short url route\n")
}