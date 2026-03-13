package handlers

import (
	"net/http"
	"time"

	"github.com/your-org/error-simulator/models"
	"github.com/your-org/error-simulator/pipeline"
)

// MultiFileOrder handles GET /error/multi-file/order.
// Call chain: handler → pipeline.ProcessOrder → formatter.BuildInvoice (panic in formatter).
func MultiFileOrder(pl *pipeline.Pipeline) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		order := &models.Order{
			ID:              "mf-ord-1",
			Amount:          49.99,
			ShippingAddress: nil, // deliberate: formatter will panic (multi-file)
			CreatedAt:       time.Now(),
		}
		_ = pl.ProcessOrder(order)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
