package handlers

import (
	"encoding/json"
	"net/http"
)

// OrderService simulates a service that processes orders. It may depend on
// other components such as a Repo. For safety, methods validate critical
// dependencies before use.
type OrderService struct {
	Repo interface{
		Save(order map[string]interface{}) error
	}
}

// ProcessOrder handles an HTTP request to process an order. It validates
// that the receiver and its dependencies are non-nil before use and returns
// appropriate HTTP errors instead of panicking.
func (s *OrderService) ProcessOrder(w http.ResponseWriter, r *http.Request) {
	if s == nil {
		http.Error(w, "internal server error: service not initialized", http.StatusInternalServerError)
		return
	}
	if s.Repo == nil {
		http.Error(w, "internal server error: repository not initialized", http.StatusInternalServerError)
		return
	}

	var order map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "bad request: invalid JSON", http.StatusBadRequest)
		return
	}

	if err := s.Repo.Save(order); err != nil {
		http.Error(w, "failed to save order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("order processed"))
}

// NilPointer is a handler factory that returns an http.HandlerFunc which
// invokes ProcessOrder on the provided service. It guards against a nil
// service by returning a handler that responds with an error.
func NilPointer(s *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s == nil {
			http.Error(w, "internal server error: service not provided", http.StatusInternalServerError)
			return
		}
		s.ProcessOrder(w, r)
	}
}
