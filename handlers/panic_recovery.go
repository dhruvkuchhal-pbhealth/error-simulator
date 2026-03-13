package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

// PaymentService is a simple example service.
type PaymentService struct{}

// PaymentRequest represents a payment request payload.
type PaymentRequest struct {
	Amount float64 `json:"amount"`
	Method string  `json:"method"`
}

// ProcessPayment processes a payment request from an http.Request.
func (p *PaymentService) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// Defensive checks to avoid panics from nil inputs.
	if p == nil {
		// Shouldn't happen in normal usage, but guard anyway.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if r == nil {
		http.Error(w, "bad request: nil request", http.StatusBadRequest)
		return
	}

	// Ensure body is not nil
	if r.Body == nil {
		http.Error(w, "bad request: empty body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Amount <= 0 {
		http.Error(w, "bad request: amount must be > 0", http.StatusBadRequest)
		return
	}
	if req.Method == "" {
		http.Error(w, "bad request: method required", http.StatusBadRequest)
		return
	}

	// Simulate processing; ensure we don't panic on unexpected internal state.
	if err := p.doProcess(&req); err != nil {
		if errors.Is(err, errUnsupportedMethod) {
			http.Error(w, "unsupported payment method", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"processed"}`))
}

var errUnsupportedMethod = errors.New("unsupported method")

// doProcess performs the internal processing and returns errors instead of panicking.
func (p *PaymentService) doProcess(req *PaymentRequest) error {
	if req == nil {
		return errors.New("nil request")
	}

	// Example: only "card" and "ach" are supported
	switch req.Method {
	case "card", "ach":
		// pretend to process
		return nil
	default:
		return errUnsupportedMethod
	}
}
