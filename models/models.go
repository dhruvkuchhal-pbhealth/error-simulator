package models

import "time"

// Order represents an order with optional patient billing info.
type Order struct {
	ID        string
	Amount    float64
	Patient   *Patient
	CreatedAt time.Time
}

// Patient holds customer/patient details with optional address.
type Patient struct {
	Name    string
	Age     int
	Address *Address
}

// Address is a nested address struct.
type Address struct {
	Street string
	City   string
	Zip    string
}

// Product represents a product in inventory or reports.
type Product struct {
	ID    string
	Name  string
	Price float64
	Stock int
}

// User represents a user row from the database.
type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
}

// ErrorEvent is the Kafka event schema consumed by ai-debugger.
type ErrorEvent struct {
	Service      string `json:"service"`
	Repository   string `json:"repository"`
	Branch       string `json:"branch"`
	ErrorMessage string `json:"error_message"`
	StackTrace   string `json:"stack_trace"`
	Timestamp    string `json:"timestamp"`
	Environment  string `json:"environment"`
}
