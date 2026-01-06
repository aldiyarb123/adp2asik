package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Server holds shared state
type Server struct {
	mu       sync.Mutex
	data     map[string]string
	requests int
}

// Constructor
func NewServer() *Server {
	return &Server{
		data: make(map[string]string),
	}
}

// Middleware to count requests
func (s *Server) countRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		s.requests++
		s.mu.Unlock()
		next(w, r)
	}
}

// POST /data
func (s *Server) postDataHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]string

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	for k, v := range body {
		s.data[k] = v
	}
	s.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

// GET /data
func (s *Server) getDataHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.data)
}

// DELETE /data/{key}
func (s *Server) deleteDataHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	delete(s.data, key)
	w.WriteHeader(http.StatusNoContent)
}

// GET /stats
func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	stats := map[string]int{
		"total_requests": s.requests,
		"database_size":  len(s.data),
	}
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Background worker
func (s *Server) startBackgroundWorker(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			fmt.Printf(
				"Status: Requests=%d, DB size=%d\n",
				s.requests,
				len(s.data),
			)
			s.mu.Unlock()
		case <-ctx.Done():
			fmt.Println("Background worker stopped")
			return
		}
	}
}

func main() {
	server := NewServer()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /data", server.countRequests(server.postDataHandler))
	mux.HandleFunc("GET /data", server.countRequests(server.getDataHandler))
	mux.HandleFunc("DELETE /data/{key}", server.countRequests(server.deleteDataHandler))
	mux.HandleFunc("GET /stats", server.countRequests(server.statsHandler))

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go server.startBackgroundWorker(ctx)

	go func() {
		fmt.Println("Server running on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done() // wait for Ctrl+C
	fmt.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpServer.Shutdown(shutdownCtx)
	fmt.Println("Server stopped gracefully")
}
