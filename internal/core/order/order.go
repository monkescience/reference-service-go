package order

import (
	"crypto/rand"
	"time"

	"github.com/oapi-codegen/runtime/types"
	"github.com/oklog/ulid/v2"
)

// OrderItem represents an item in an order
type OrderItem struct {
	Name string
}

// Order represents an order in the system
type Order struct {
	OrderID      string
	CustomerID   types.UUID
	CreationDate time.Time
	Status       string
	Items        []OrderItem
}

// OrderRequest represents a request to create an order
type OrderRequest struct {
	CustomerID types.UUID
	Items      []OrderItem
}

// Service handles the business logic for orders
type Service struct {
	repository Repository
}

// Repository defines the interface for order storage
type Repository interface {
	StoreOrder(order Order) error
	GetOrder(orderID string) (Order, bool)
	GetOrders(customerID *string, limit, offset int) []Order
}

// NewService creates a new order service
func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// CreateOrder creates a new order with the given details
func (s *Service) CreateOrder(orderRequest OrderRequest) (Order, error) {
	// Create the order with a ULID-based ID
	order := Order{
		OrderID:      GenerateOrderID(),
		CustomerID:   orderRequest.CustomerID,
		CreationDate: time.Now().UTC(),
		Status:       "order_placed",
		Items:        orderRequest.Items,
	}

	// Store the order
	if err := s.repository.StoreOrder(order); err != nil {
		return Order{}, err
	}

	return order, nil
}

// GetOrder retrieves an order by ID
func (s *Service) GetOrder(orderID string) (Order, bool) {
	return s.repository.GetOrder(orderID)
}

// GetOrders retrieves orders with optional filtering and pagination
func (s *Service) GetOrders(customerID *string, limit, offset int) []Order {
	return s.repository.GetOrders(customerID, limit, offset)
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
