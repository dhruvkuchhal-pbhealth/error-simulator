package handlers

import (
	"errors"
	"fmt"
	"net/http"
)

// PaymentService simulates a service that processes payments.
type PaymentService struct {
	MaxAmount int64
}

// ProcessPayment processes the payment of the given amount.
// It returns an error instead of panicking when validation fails.
func (p *PaymentService) ProcessPayment(amount int64) error {
	if amount <= 0 {
		return errors.New("payment amount must be positive")
	}

	if amount > p.MaxAmount {
		// Return a descriptive error instead of panicking.
		return fmt.Errorf("payment amount exceeds maximum transaction limit: got %d, max %d", amount, p.MaxAmount)
	}

	// Simulate processing...
	_ = amount // placeholder for real processing logic
	return nil
}

// Handler wrapper that demonstrates converting service errors to HTTP responses.
func PaymentHandler(svc *PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For the sake of example, amount is read from query param "amount" as int64.
		// In real code, parse body or form values properly and handle errors.
		amountStr := r.URL.Query().Get("amount")
		var amount int64
		if amountStr == "" {
			http.Error(w, "missing amount", http.StatusBadRequest)
			return
		}
		_, err := fmt.Sscan(amountStr, &amount)
		if err != nil {
			http.Error(w, "invalid amount", http.StatusBadRequest)
			return
		}

		if err := svc.ProcessPayment(amount); err != nil {
			// Convert validation error to 400 Bad Request
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("payment processed"))
	}
}
