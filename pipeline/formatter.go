package pipeline

import "github.com/your-org/error-simulator/models"

// BuildInvoice builds an invoice line for an order including shipping address.
// BUG: No nil check on order.ShippingAddress — panic when ShippingAddress is nil (multi-file stack).
func BuildInvoice(order *models.Order) string {
	if order == nil {
		return ""
	}
	// Intentionally dereference nil ShippingAddress so panic happens in this file (layer 3).
	street := order.ShippingAddress.Street
	city := order.ShippingAddress.City
	return street + ", " + city
}
