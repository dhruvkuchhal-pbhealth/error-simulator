package pipeline

import (
	"errors"
	"fmt"
)

// Invoice represents a generated invoice for an order.
type Invoice struct {
	OrderID   string
	LineItems []LineItem
	Total     float64
}

// LineItem represents a single line in the invoice.
type LineItem struct {
	SKU      string
	Quantity int
	Price    float64
}

// Order represents an order to be processed into an invoice.
type Order struct {
	ID        string
	Items     []OrderItem
	Customer  *Customer
	ShipState string
}

// OrderItem represents an item in an order.
type OrderItem struct {
	SKU      string
	Quantity int
	Price    float64
}

// Customer represents customer details.
type Customer struct {
	Name string
	ID   string
}

// BuildInvoice constructs an Invoice from an Order. Returns an error when the order is nil
// or contains invalid data instead of panicking.
func BuildInvoice(o *Order) (*Invoice, error) {
	if o == nil {
		return nil, errors.New("order is nil")
	}

	if o.ID == "" {
		return nil, errors.New("order ID is empty")
	}

	inv := &Invoice{
		OrderID:   o.ID,
		LineItems: make([]LineItem, 0, len(o.Items)),
	}

	var total float64
	for i := range o.Items {
		it := &o.Items[i]
		// defensive checks for each item
		if it == nil { // can't actually be nil since slice of structs, but keep check if types change
			return nil, fmt.Errorf("order %s: item at index %d is nil", o.ID, i)
		}
		if it.Quantity <= 0 {
			return nil, fmt.Errorf("order %s: invalid quantity for sku %s", o.ID, it.SKU)
		}
		if it.Price < 0 {
			return nil, fmt.Errorf("order %s: negative price for sku %s", o.ID, it.SKU)
		}

		li := LineItem{
			SKU:      it.SKU,
			Quantity: it.Quantity,
			Price:    it.Price,
		}
		inv.LineItems = append(inv.LineItems, li)
		total += float64(it.Quantity) * it.Price
	}

	inv.Total = total
	return inv, nil
}
