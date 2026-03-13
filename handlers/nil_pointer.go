package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/your-org/error-simulator/models"
)

// OrderService processes orders and generates invoices.
// The bug: ProcessOrder does not nil-check order.Patient before accessing nested fields.
type OrderService struct {
	baseURL string
}

// NewOrderService returns an OrderService with default config.
func NewOrderService() *OrderService {
	return &OrderService{baseURL: "https://api.example.com"}
}

// ProcessOrder validates the order, resolves patient billing info, and builds an invoice line.
// When order.Patient is nil, accessing order.Patient.Name causes a nil pointer dereference.
func (s *OrderService) ProcessOrder(order *models.Order) (invoiceLine string, err error) {
	if order == nil {
		return "", nil
	}
	// Simulate fetching tax rate and formatting amount.
	amountStr := formatCurrency(order.Amount)
	created := order.CreatedAt.Format(time.RFC3339)
	// Build billing line: we need patient name and address for the invoice.
	// BUG: No nil check on order.Patient — if the order was created without
	// patient data (e.g. anonymous checkout), Patient is nil.
	patientName := order.Patient.Name
	patientCity := order.Patient.Address.City
	invoiceLine = patientName + " | " + patientCity + " | " + amountStr + " | " + created
	return invoiceLine, nil
}

func formatCurrency(amount float64) string {
	return fmt.Sprintf("USD %.2f", amount)
}

// NilPointer handles GET /error/nil-pointer.
// It creates an Order with nil Patient and calls ProcessOrder to trigger the dereference.
func NilPointer(svc *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		order := &models.Order{
			ID:        "ord-12345",
			Amount:    99.99,
			Patient:   nil, // deliberately nil to trigger bug
			CreatedAt: time.Now(),
		}
		_, _ = svc.ProcessOrder(order)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
