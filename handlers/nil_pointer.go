package handlers

import (
	"encoding/json"
	"net/http"
)

// OrderService represents a service that processes orders.
// Some internals may be nil in tests or misconfiguration; guard against that.
type OrderService struct {
	// Example internal pointer (could be anything that might be nil)
	Validator func(interface{}) bool
}

// ProcessOrder processes an order payload. It defensively checks for nil receiver
// and nil internal fields to avoid panics.
func (s *OrderService) ProcessOrder(r *http.Request) (int, interface{}, error) {
	// Guard against nil receiver
	if s == nil {
		return http.StatusInternalServerError, map[string]string{"error": "order service unavailable"}, nil
	}

	// Read body into a generic map
	var payload interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, map[string]string{"error": "invalid payload"}, err
	}

	// If a Validator is provided, use it; otherwise accept.
	if s.Validator != nil {
		if ok := s.Validator(payload); !ok {
			return http.StatusBadRequest, map[string]string{"error": "validation failed"}, nil
		}
	}

	// Simulate processing and return success
	return http.StatusOK, map[string]string{"status": "processed"}, nil
}

// NilPointer returns an http.HandlerFunc that ensures a non-nil OrderService is used.
func NilPointer(svc *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If svc is nil, construct a safe default to avoid dereference
		if svc == nil {
			svc = &OrderService{}
		}

		status, resp, err := svc.ProcessOrder(r)
		if err != nil {
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
