package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tinylynx/internal/concurrency"
	"tinylynx/internal/handler"
	"tinylynx/internal/storage"
)

func main() {
	// Run database migrations
	storage.RunMigrations()

	// Initialize worker pool for analytics
	workerPool := concurrency.NewWorkerPool(10)
	workerPool.Start()
	defer workerPool.Shutdown()

	// Initialize rate limiter
	rateLimiter := concurrency.NewRateLimiter(100, 1*time.Minute)

	// Initialize handlers
	linkHandler := handler.NewLinkHandler()

	// Setup routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow(r.RemoteAddr) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		linkHandler.CreateShortLink(w, r)
	})

	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow(r.RemoteAddr) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		linkHandler.GetLinkStats(w, r)
	})

	mux.HandleFunc("/api/analytics", func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow(r.RemoteAddr) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		linkHandler.GetLinkAnalytics(w, r)
	})

	// Redirect routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow(r.RemoteAddr) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		linkHandler.Redirect(w, r)
	})

	// Get server port from env or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("TinyLynx server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}