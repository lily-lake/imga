package api

import (
	// "encoding/json"
	"fmt"
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
	fmt.Fprintf(w, "create short url route\n")
}