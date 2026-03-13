package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const maxTransactionLimit = 10000.0

// PaymentService processes payments. The bug: ProcessPayment panics when
// amount exceeds the limit instead of returning an error.
type PaymentService struct {
	merchantID string
}

// NewPaymentService returns a payment service with default merchant.
func NewPaymentService() *PaymentService {
	return &PaymentService{merchantID: "merchant_default"}
}

// ProcessPayment validates and processes a payment. In production this would
// call a payment gateway. Here we explicitly panic when amount > maxTransactionLimit
// to simulate a business rule enforced via panic.
func (s *PaymentService) ProcessPayment(amount float64, currency string) (txID string, err error) {
		log.Printf("[PaymentService] ProcessPayment started merchant=%s amount=%.2f currency=%s",
			s.merchantID, amount, currency)
		// No validation — direct panic when amount exceeds limit.
	if amount > maxTransactionLimit {
		panic(fmt.Sprintf("payment amount exceeds maximum transaction limit: got %v, max %v",
			amount, maxTransactionLimit))
	}
	txID = fmt.Sprintf("tx_%d", time.Now().UnixNano())
	return txID, nil
}

// PanicRecovery handles GET /error/panic.
// It calls ProcessPayment with an amount over the limit to trigger the panic.
func PanicRecovery(svc *PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txID, err := svc.ProcessPayment(999999, "USD")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(txID))
	}
}
