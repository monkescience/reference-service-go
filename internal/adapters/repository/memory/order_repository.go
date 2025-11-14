package memory

import (
	"sync"

	dom "reference-service-go/internal/domain/order"
	"reference-service-go/internal/ports/order"
)

// Ensure repository implements the port
var _ order.Repository = (*Repository)(nil)

// Repository is an in-memory implementation of order.Repository
type Repository struct {
	orders      map[string]dom.Order
	ordersMutex sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{orders: make(map[string]dom.Order)}
}

func (r *Repository) StoreOrder(o dom.Order) error {
	r.ordersMutex.Lock()
	defer r.ordersMutex.Unlock()
	r.orders[o.OrderID] = o
	return nil
}

func (r *Repository) GetOrder(orderID string) (dom.Order, bool) {
	r.ordersMutex.RLock()
	defer r.ordersMutex.RUnlock()
	o, exists := r.orders[orderID]
	return o, exists
}

func (r *Repository) GetOrders(customerID *string, limit, offset int) []dom.Order {
	r.ordersMutex.RLock()
	defer r.ordersMutex.RUnlock()

	var filtered []dom.Order
	for _, o := range r.orders {
		if customerID != nil && o.CustomerID != *customerID {
			continue
		}
		filtered = append(filtered, o)
	}

	var result []dom.Order
	for i, o := range filtered {
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
