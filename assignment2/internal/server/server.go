package server

import (
	"assignment2/internal/storage"
	"sync"
	"time"
)

type Server struct {
	store     *storage.MemoryStore
	mu        sync.Mutex
	requests  int
	startTime time.Time
}

func NewServer() *Server {
	return &Server{
		store:     storage.NewMemoryStore(),
		startTime: time.Now(),
	}
}

func (s *Server) IncrementRequests() {
	s.mu.Lock()
	s.requests++
	s.mu.Unlock()
}

func (s *Server) Stats() (int, int, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.requests, s.store.Size(), int(time.Since(s.startTime).Seconds())
}
