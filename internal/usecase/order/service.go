package order

import (
	"time"

	domain "reference-service-go/internal/domain/order"
	"reference-service-go/internal/ports/order"
)

// Service handles the application use cases for orders
// Depends on a port.Repository and orchestrates domain actions
type Service struct {
	repo order.Repository
}

// Ensure Service implements the input port
var _ order.Service = (*Service)(nil)

// NewService creates a new order service
func NewService(repo order.Repository) *Service {
	return &Service{repo: repo}
}

// CreateOrder creates a new order with the given details
func (s *Service) CreateOrder(request domain.OrderRequest) (domain.Order, error) {
	id, err := domain.GenerateOrderID(request.Country)
	if err != nil {
		return domain.Order{}, err
	}
	ord := domain.Order{
		OrderID:      id,
		CustomerID:   request.CustomerID,
		CreationDate: time.Now().UTC(),
		Status:       "order_placed",
		Items:        request.Items,
	}
	if err := s.repo.StoreOrder(ord); err != nil {
		return domain.Order{}, err
	}
	return ord, nil
}

// GetOrder retrieves an order by ID
func (s *Service) GetOrder(orderID string) (domain.Order, bool) {
	return s.repo.GetOrder(orderID)
}

// GetOrders retrieves orders with optional filtering and pagination
func (s *Service) GetOrders(customerID *string, limit, offset int) []domain.Order {
	return s.repo.GetOrders(customerID, limit, offset)
}
