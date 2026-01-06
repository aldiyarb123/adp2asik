package server

import (
	"context"
	"log"
	"time"
)

func (s *Server) StartWorker(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			req, size, _ := s.Stats()
			log.Printf("[WORKER] requests=%d db_size=%d\n", req, size)

		case <-ctx.Done():
			log.Println("[WORKER] stopped")
			return
		}
	}
}
