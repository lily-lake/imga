package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"sync"
	"testing"
)

var (
	mockURLMap   = make(map[string]string)
	mockURLMapMu sync.RWMutex
)

func TestGenerateShortCode(t *testing.T) {
	t.Run("Generates correct length short code", func(t *testing.T) {

		code := generateShortCode()
		if len(code) != 6 {
			t.Errorf("expected short code length of 6, got %d", len(code))
		}
	})
}

func TestGetUniqueShortCode(t *testing.T) {
	t.Run("Generates unique short code", func(t *testing.T) {
		// I know this is a silly test, would spend more time in production code :)
		code1 := getUniqueShortCode(mockURLMap, &mockURLMapMu)
		code2 := getUniqueShortCode(mockURLMap, &mockURLMapMu)
		if code1 == code2 {
			t.Errorf("expected unique short codes, got %s and %s", code1, code2)
		}
	})
}

func TestIsValidURL(t *testing.T) {
	t.Run("Returns true for valid URL", func(t *testing.T) {
		url := "https://www.google.com"
		if !isValidURL(url) {
			t.Errorf("expected valid URL, got %s", url)
		}
	})
	t.Run("Returns false for invalid URL", func(t *testing.T) {
		url := "invalid-url"
		if isValidURL(url) {
			t.Errorf("expected invalid URL, got %s", url)
		}
	})
}

func TestRedirectToOriginalURLHandler(t *testing.T) {
	mockURLMap["111111"] = "https://www.example.com"
	t.Run("Redirects to original URL", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/111111", nil)
		response := httptest.NewRecorder()
		RedirectToOriginalURLHandler(mockURLMap, &mockURLMapMu)(response, request)
		if response.Code != 302 {
			t.Errorf("expected 302 status code, got %d", response.Code)
		}
		if response.Header().Get("Location") != "https://www.example.com" {
			t.Errorf("expected Location header to be https://www.example.com, got %s", response.Header().Get("Location"))
		}
	})
}

func TestCreateShortURLHandler(t *testing.T) {
	t.Run("Creates short URL", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(`{"url": "https://www.example.com"}`))
		response := httptest.NewRecorder()
		CreateShortURLHandler(mockURLMap, &mockURLMapMu)(response, request)
		if response.Code != 200 {
			t.Errorf("expected 200 status code, got %d", response.Code)
		}
		var responseBody ShortenURLResponse
		if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
			t.Errorf("expected valid JSON response, got %s", response.Body.String())
		}
		if responseBody.OriginalURL != "https://www.example.com" {
			t.Errorf("expected OriginalURL to be https://www.example.com, got %s", responseBody.OriginalURL)
		}
		// I believe different handler structure would be needed to test that the short code generated is the one in the response
		// I don't think this would be a high value test
	})
	t.Run("Creates short URL with custom short code", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(`{"url": "https://www.example.com", "shortcode": "abcdefg"}`))
		response := httptest.NewRecorder()
		CreateShortURLHandler(mockURLMap, &mockURLMapMu)(response, request)
		if response.Code != 200 {
			t.Errorf("expected 200 status code, got %d", response.Code)
		}
		var responseBody ShortenURLResponse
		if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
			t.Errorf("expected valid JSON response, got %s", response.Body.String())
		}
		if responseBody.OriginalURL != "https://www.example.com" {
			t.Errorf("expected OriginalURL to be https://www.example.com, got %s", responseBody.OriginalURL)
		}
	})
	t.Run("Returns error if custom short code already exists", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(`{"url": "https://www.example.com", "shortcode": "abcdefg"}`))
		response := httptest.NewRecorder()
		CreateShortURLHandler(mockURLMap, &mockURLMapMu)(response, request)
		if response.Code != 409 {
			t.Errorf("expected 409 status code, got %d", response.Code)
		}
	})
}
