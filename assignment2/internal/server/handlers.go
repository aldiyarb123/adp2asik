package server

import (
	"encoding/json"
	"net/http"
)

// POST /data
func (s *Server) PostData(w http.ResponseWriter, r *http.Request) {
	s.IncrementRequests()

	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for k, v := range payload {
		s.store.Set(k, v)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "stored"})
}

// GET /data
func (s *Server) GetData(w http.ResponseWriter, r *http.Request) {
	s.IncrementRequests()
	json.NewEncoder(w).Encode(s.store.GetAll())
}

// DELETE /data/{key}
func (s *Server) DeleteData(w http.ResponseWriter, r *http.Request) {
	s.IncrementRequests()

	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "Key required", http.StatusBadRequest)
		return
	}

	if !s.store.Delete(key) {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"deleted": key})
}

// GET /stats
func (s *Server) StatsHandler(w http.ResponseWriter, r *http.Request) {
	s.IncrementRequests()

	req, size, uptime := s.Stats()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_requests": req,
		"database_size":  size,
		"uptime_seconds": uptime,
	})
}
