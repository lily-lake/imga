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
	"sync"
)

// Shorten URL params
type ShortenURLParams struct {
	URL       string  `json:"URL"`
	ShortCode *string `json:"ShortCode,omitempty"` // Optional custom short code
}

// Shorten URL response
type ShortenURLResponse struct {
	ShortCode   string `json:"ShortCode"`
	ShortURL    string `json:"ShortURL"`
	OriginalURL string `json:"OriginalURL"`
}

type Error struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

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
func getUniqueShortCode(urlMap map[string]string, mu *sync.RWMutex) string {
	for {
		code := generateShortCode()
		mu.RLock()
		_, exists := urlMap[code]
		mu.RUnlock()
		if !exists {
			return code
		}
		// If code exists, try again
	}
}

func CreateShortURLHandler(urlMap map[string]string, mu *sync.RWMutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var shortCode string

		// Check if custom short code is provided
		if params.ShortCode != nil && *params.ShortCode != "" {
			customCode := *params.ShortCode
			// Check if custom code already exists
			mu.RLock()
			_, exists := urlMap[customCode]
			mu.RUnlock()
			if exists {
				http.Error(w, "Short code already exists", http.StatusConflict)
				log.Printf("Custom short code already exists: %s", customCode)
				return
			}
			shortCode = customCode
		} else {
			// Generate unique short code if no custom code provided
			shortCode = getUniqueShortCode(urlMap, mu)
		}

		// Store the short code
		mu.Lock()
		urlMap[shortCode] = params.URL
		mu.Unlock()

		// Create response
		response := ShortenURLResponse{
			ShortCode:   shortCode,
			ShortURL:    fmt.Sprintf("http://localhost:8080/%s", shortCode),
			OriginalURL: params.URL,
		}

		log.Printf("shorten request for URL: %s", params.URL)

		// Set response headers and status code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Encode response as JSON
		json.NewEncoder(w).Encode(response)
	}
}

func RedirectToOriginalURLHandler(urlMap map[string]string, mu *sync.RWMutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Println("invalid request type to /{shortCode}")
			return
		}

		// Get short code from request
		shortCode := strings.TrimPrefix(r.URL.Path, "/")

		// Get original URL from map
		mu.RLock()
		originalURL, exists := urlMap[shortCode]
		mu.RUnlock()
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
}
