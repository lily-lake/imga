package api

import (
	"testing"
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
		code1 := getUniqueShortCode()
		code2 := getUniqueShortCode()
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
