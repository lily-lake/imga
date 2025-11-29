package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"imga/api"
)

func main() {
	urlMap := make(map[string]string)
	var urlMapMu sync.RWMutex

	http.HandleFunc("/shorten", api.CreateShortURLHandler(urlMap, &urlMapMu))
	http.HandleFunc("/{shortCode}", api.RedirectToOriginalURLHandler(urlMap, &urlMapMu))

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil, // Uses http.DefaultServeMux
	}

	// Start server in a goroutine
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
