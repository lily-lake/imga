package api

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Shorten URL params
type ShortenURLParams struct {
	URL string
}

// Shorten URL response
type ShortenURLResponse struct {
	ShortCode   string
	ShortURL    string
	OriginalURL string
}

type Error struct {
	Code    int
	Message string
}

// URL short code map
var (
	URLMap = make(map[string]string)
)

// URL validator
func isValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

// Generate short code
func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		// Generate random byte
		randomByte := make([]byte, 1)
		rand.Read(randomByte)
		b[i] = charset[randomByte[0]%byte(len(charset))]
	}
	return string(b)
}

// Generate code that is not already in the map
func getUniqueShortCode() string {
	for {
		code := generateShortCode()
		if _, exists := URLMap[code]; !exists {
			return code
		}
		// If code exists, try again
	}
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

	// Generate unique short code
	shortCode := getUniqueShortCode()
	URLMap[shortCode] = params.URL

	// Create response
	response := ShortenURLResponse{
		ShortCode:   shortCode,
		ShortURL:    fmt.Sprintf("http://localhost:8080/%s", shortCode),
		OriginalURL: params.URL,
	}

	log.Printf("shorten request for URL: %s", params.URL)

	// Encode response as JSON
	json.NewEncoder(w).Encode(response)
}

func RedirectToOriginalURLHandler(w http.ResponseWriter, r *http.Request) {
	// Get short code from request
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	// Get original URL from map
	originalURL, exists := URLMap[shortCode]
	if !exists {
		http.Error(w, "Short code not found", http.StatusNotFound)
		log.Printf("Short code not found: %s", shortCode)
		return
	}

	log.Printf("Redirecting to original URL: %s", originalURL)

	// Redirect to original URL
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusFound)
}
