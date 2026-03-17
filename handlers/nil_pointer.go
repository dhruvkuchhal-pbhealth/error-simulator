package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

// OrderService simulates a service that processes orders.
// Fields may be nil in some situations in the original code; ensure safe access.
type OrderService struct {
	// Imagine some internal client or dependency that might be nil
	Client interface{}
}

// ProcessOrder processes an order; returns an error instead of panicking on nil receiver/deps.
func (s *OrderService) ProcessOrder() error {
	if s == nil {
		return errors.New("order service is not initialized")
	}
	// If the service depends on Client, ensure it's non-nil before use.
	if s.Client == nil {
		// Either initialize default behavior or return an explicit error.
		return errors.New("order service client is nil")
	}

	// Normal processing would occur here. Keep minimal to avoid panic.
	return nil
}

// NilPointer returns an HTTP handler that uses OrderService safely.
func NilPointer(svc *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate service before calling methods.
		if svc == nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "order service not available"})
			return
		}

		if err := svc.ProcessOrder(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "processed"})
	}
}
