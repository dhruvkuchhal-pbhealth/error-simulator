package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

// PaymentService is a simple example service that processes payments.
// In the original code it panicked on invalid input; we return errors instead.
type PaymentService struct{
	MaxAmount int64
}

// PaymentRequest represents a payment request payload.
type PaymentRequest struct{
	AccountID string `json:"account_id"`
	Amount    int64  `json:"amount"`
}

// ProcessPayment validates and processes a payment request.
// It returns an error for invalid input instead of panicking.
func (s *PaymentService) ProcessPayment(r *http.Request) error {
	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// use configured MaxAmount if set, otherwise default to 10000
	max := s.MaxAmount
	if max == 0 {
		max = 10000
	}

	if req.Amount > max {
		return errors.New("payment amount exceeds maximum transaction limit")
	}

	// Placeholder for actual processing logic. Return nil to indicate success.
	return nil
}

// Handler to expose ProcessPayment over HTTP
func (s *PaymentService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := s.ProcessPayment(r); err != nil {
		// map validation error to 400 Bad Request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
