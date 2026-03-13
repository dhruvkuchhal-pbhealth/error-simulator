package handlers

import (
	"encoding/json"
	"net/http"
)

// OrderService simulates an order processing service.
type OrderService struct {
	// maybe some dependencies would go here
	Name string
}

// ProcessOrder processes an order. It previously assumed the receiver was non-nil and dereferenced it,
// causing a panic when called on a nil *OrderService. We now guard against a nil receiver and return
// an error response instead of panicking.
func (s *OrderService) ProcessOrder(w http.ResponseWriter) {
	if s == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "order service not initialized"})
		return
	}

	// normal processing (kept minimal as original intent was demonstration)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "processed", "service": s.Name})
}

// NilPointer returns an HTTP handler that demonstrates calling ProcessOrder.
// The handler also guards against a nil service before calling the method.
func NilPointer(s *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s == nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "order service not available"})
			return
		}
		s.ProcessOrder(w)
	}
}
