package order

import (
	"sync"

	"reference-service-go/internal/core/order"
)

// Repository implements the core.Repository interface for storing orders in memory
type Repository struct {
	orders      map[string]order.Order
	ordersMutex sync.RWMutex
}

// NewRepository creates a new order repository
func NewRepository() *Repository {
	return &Repository{
		orders: make(map[string]order.Order),
	}
}

// StoreOrder stores an order in the repository
func (r *Repository) StoreOrder(o order.Order) error {
	r.ordersMutex.Lock()
	defer r.ordersMutex.Unlock()

	r.orders[o.OrderID] = o
	return nil
}

// GetOrder retrieves an order by ID
func (r *Repository) GetOrder(orderID string) (order.Order, bool) {
	r.ordersMutex.RLock()
	defer r.ordersMutex.RUnlock()

	o, exists := r.orders[orderID]
	return o, exists
}

// GetOrders retrieves orders with optional filtering and pagination
func (r *Repository) GetOrders(customerID *string, limit, offset int) []order.Order {
	r.ordersMutex.RLock()
	defer r.ordersMutex.RUnlock()

	// Filter orders by customer_id if provided
	var filteredOrders []order.Order
	for _, o := range r.orders {
		if customerID != nil && o.CustomerID.String() != *customerID {
			continue
		}
		filteredOrders = append(filteredOrders, o)
	}

	// Apply pagination
	var result []order.Order
	for i, o := range filteredOrders {
		if i < offset {
			continue
		}
		if len(result) >= limit {
			break
		}
		result = append(result, o)
	}

	return result
}
