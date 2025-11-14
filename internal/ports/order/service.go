// Package order defines inbound ports for the order use case layer.
package order

import dom "reference-service-go/internal/domain/order"

// Service is the input port for order-related use cases.
// Implemented by internal/usecase/order.Service and consumed by adapters.
type Service interface {
	CreateOrder(request dom.OrderRequest) (dom.Order, error)
	GetOrder(orderID string) (dom.Order, bool)
	GetOrders(customerID *string, limit, offset int) []dom.Order
}
