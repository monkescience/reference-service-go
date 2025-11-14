package memory

import (
	"sort"
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

func (repository *Repository) StoreOrder(ord dom.Order) error {
	repository.ordersMutex.Lock()
	defer repository.ordersMutex.Unlock()
	repository.orders[ord.OrderID] = ord
	return nil
}

func (repository *Repository) GetOrder(orderID string) (dom.Order, bool) {
	repository.ordersMutex.RLock()
	defer repository.ordersMutex.RUnlock()
	ord, exists := repository.orders[orderID]
	return ord, exists
}

func (repository *Repository) GetOrders(customerID *string, limit, offset int) []dom.Order {
	repository.ordersMutex.RLock()
	defer repository.ordersMutex.RUnlock()

	var filteredOrders []dom.Order
	for _, ord := range repository.orders {
		if customerID != nil && ord.CustomerID != *customerID {
			continue
		}
		filteredOrders = append(filteredOrders, ord)
	}

	// Sort deterministically: by CreationDate asc, then OrderID asc
	sort.Slice(filteredOrders, func(i, j int) bool {
		if filteredOrders[i].CreationDate.Equal(filteredOrders[j].CreationDate) {
			return filteredOrders[i].OrderID < filteredOrders[j].OrderID
		}
		return filteredOrders[i].CreationDate.Before(filteredOrders[j].CreationDate)
	})

	var result []dom.Order
	for index, ord := range filteredOrders {
		if index < offset {
			continue
		}
		if len(result) >= limit {
			break
		}
		result = append(result, ord)
	}
	return result
}
