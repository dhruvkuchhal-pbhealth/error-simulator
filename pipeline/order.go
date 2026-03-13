package pipeline

import "github.com/your-org/error-simulator/models"

// Pipeline processes orders through the formatter (layer 2). Panic occurs in formatter (layer 3).
type Pipeline struct{}

// ProcessOrder validates the order and builds an invoice via the formatter.
// Panic occurs in formatter.BuildInvoice when order.ShippingAddress is nil.
func (p *Pipeline) ProcessOrder(order *models.Order) (invoice string) {
	if order == nil {
		return ""
	}
	return BuildInvoice(order)
}
