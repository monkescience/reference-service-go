package order

import dom "reference-service-go/internal/domain/order"

// Repository defines the interface for order storage ports
// Adapters (e.g., memory, db) must implement this.
type Repository interface {
	StoreOrder(order dom.Order) error
	GetOrder(orderID string) (dom.Order, bool)
	GetOrders(customerID *string, limit, offset int) []dom.Order
}
