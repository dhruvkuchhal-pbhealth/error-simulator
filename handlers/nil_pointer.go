package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

// OrderService simulates an external service used by the handler.
// In real code this might have clients or other dependencies that may be nil.
type OrderService struct {
	// Example dependency that could be nil
	Publisher *KafkaPublisher
}

// KafkaPublisher is a placeholder for an external publisher that may be nil.
type KafkaPublisher struct {
	Topic string
}

// ProcessOrder processes an order. It returns an error instead of panicking when
// required dependencies are nil.
func (s *OrderService) ProcessOrder() error {
	if s == nil {
		return errors.New("order service is nil")
	}
	if s.Publisher == nil {
		return errors.New("publisher is nil")
	}
	// Simulate publishing; in real code this would publish to kafka etc.
	// If Publisher.Topic is empty, treat as error
	if s.Publisher.Topic == "" {
		return errors.New("publisher topic is empty")
	}
	// Success
	return nil
}

// NilPointer returns an HTTP handler that uses OrderService.
// The handler validates the service and returns proper HTTP errors instead of panicking.
func NilPointer(svc *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate service before use
		if svc == nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "order service unavailable"})
			return
		}

		if err := svc.ProcessOrder(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "processed"})
	}
}
