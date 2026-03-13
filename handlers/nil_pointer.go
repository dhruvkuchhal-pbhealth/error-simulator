package handlers

import (
	"encoding/json"
	"net/http"
)

// Order represents an order payload.
type Order struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

// Validator is a simple interface to validate orders.
type Validator interface {
	Validate(o *Order) error
}

// OrderService handles order processing.
type OrderService struct {
	validator Validator
}

// ProcessOrder processes an order. It validates the order and responds accordingly.
func (s *OrderService) ProcessOrder(w http.ResponseWriter, r *http.Request) {
	// Defensive checks to avoid nil pointer dereference.
	if s == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if s.validator == nil {
		http.Error(w, "service misconfigured: validator is nil", http.StatusInternalServerError)
		return
	}

	var o Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}

	if err := s.validator.Validate(&o); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("order processed"))
}

// NilPointer returns an http handler that demonstrates nil pointer handling.
func NilPointer(s *OrderService) http.Handler {
	// Wrap to ensure nil service is handled gracefully by the handler itself.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s == nil {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		s.ProcessOrder(w, r)
	})
}
