package order

import (
	"time"

	dom "reference-service-go/internal/domain/order"
	"reference-service-go/internal/ports/order"
)

// Service handles the application use cases for orders
// Depends on a port.Repository and orchestrates domain actions
type Service struct {
	repo order.Repository
}

// NewService creates a new order service
func NewService(repo order.Repository) *Service {
	return &Service{repo: repo}
}

// CreateOrder creates a new order with the given details
func (s *Service) CreateOrder(req dom.OrderRequest) (dom.Order, error) {
	ord := dom.Order{
		OrderID:      dom.GenerateOrderID(),
		CustomerID:   req.CustomerID,
		CreationDate: time.Now().UTC(),
		Status:       "order_placed",
		Items:        req.Items,
	}
	if err := s.repo.StoreOrder(ord); err != nil {
		return dom.Order{}, err
	}
	return ord, nil
}

// GetOrder retrieves an order by ID
func (s *Service) GetOrder(orderID string) (dom.Order, bool) {
	return s.repo.GetOrder(orderID)
}

// GetOrders retrieves orders with optional filtering and pagination
func (s *Service) GetOrders(customerID *string, limit, offset int) []dom.Order {
	return s.repo.GetOrders(customerID, limit, offset)
}
