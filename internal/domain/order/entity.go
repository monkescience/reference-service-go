package order

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// OrderItem represents an item in an order
type OrderItem struct {
	Name string
}

// Order represents an order in the system
// Note: keep domain free of transport-specific types; use primitive string for UUIDs
// Mapping from/to API types happens in adapters.
type Order struct {
	OrderID      string
	CustomerID   string
	CreationDate time.Time
	Status       string
	Items        []OrderItem
}

// OrderRequest represents a request to create an order
// Note: domain-friendly input kept simple
type OrderRequest struct {
	CustomerID string
	Items      []OrderItem
}

// GenerateOrderID generates a unique order ID using ULID
// Format: <ULID-PART-1>-NONE-<ULID-PART-2>
func GenerateOrderID() string {
	// Generate ULID
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

	// Get the ULID string
	ulidStr := id.String()

	// Split the ULID string into two parts
	// ULID is 26 characters long, so we'll split it at the middle (13 characters each)
	part1 := ulidStr[:13]
	part2 := ulidStr[13:]

	// Format according to spec: <ULID-PART-1>-NONE-<ULID-PART-2>
	return part1 + "-NONE-" + part2
}
