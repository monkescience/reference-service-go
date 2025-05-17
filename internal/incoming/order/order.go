package order

import (
	"encoding/json"
	"net/http"

	"reference-service-go/internal/core/order"
	outgoing "reference-service-go/internal/outgoing/order"
)

// Server implements the ServerInterface
type Server struct {
	service *order.Service
}

// NewServer creates a new Server instance
func NewServer() *Server {
	repository := outgoing.NewRepository()
	service := order.NewService(repository)
	return &Server{
		service: service,
	}
}

// GetOrders returns a list of orders with optional filtering
func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request, params GetOrdersParams) {
	// Default values for limit and offset
	limit := 50
	offset := 0

	if params.Limit != nil {
		limit = *params.Limit
	}

	if params.Offset != nil {
		offset = *params.Offset
	}

	// Convert UUID to string if provided
	var customerIDStr *string
	if params.CustomerId != nil {
		id := params.CustomerId.String()
		customerIDStr = &id
	}

	// Get orders from the service
	coreOrders := s.service.GetOrders(customerIDStr, limit, offset)

	// Convert core orders to API response format
	var result OrdersResponse
	for _, coreOrder := range coreOrders {
		// Convert order items
		items := make([]OrderItemResponse, len(coreOrder.Items))
		for i, item := range coreOrder.Items {
			items[i] = OrderItemResponse{
				Name: item.Name,
			}
		}

		// Create API response order
		apiOrder := OrderResponse{
			OrderId:      coreOrder.OrderID,
			CustomerId:   coreOrder.CustomerID,
			CreationDate: coreOrder.CreationDate,
			Status:       OrderStatus(coreOrder.Status),
			Items:        items,
		}

		result = append(result, apiOrder)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// PostOrders creates a new order
func (s *Server) PostOrders(w http.ResponseWriter, r *http.Request) {
	var apiOrderRequest OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&apiOrderRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert API request to core request
	coreItems := make([]order.OrderItem, len(apiOrderRequest.Items))
	for i, item := range apiOrderRequest.Items {
		coreItems[i] = order.OrderItem{
			Name: item.Name,
		}
	}

	coreOrderRequest := order.OrderRequest{
		CustomerID: apiOrderRequest.CustomerId,
		Items:      coreItems,
	}

	// Create the order using the service
	coreOrder, err := s.service.CreateOrder(coreOrderRequest)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Convert core order to API response
	items := make([]OrderItemResponse, len(coreOrder.Items))
	for i, item := range coreOrder.Items {
		items[i] = OrderItemResponse{
			Name: item.Name,
		}
	}

	apiOrder := OrderResponse{
		OrderId:      coreOrder.OrderID,
		CustomerId:   coreOrder.CustomerID,
		CreationDate: coreOrder.CreationDate,
		Status:       OrderStatus(coreOrder.Status),
		Items:        items,
	}

	// Return the created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiOrder)
}

// GetOrder returns details of a specific order
func (s *Server) GetOrder(w http.ResponseWriter, r *http.Request, orderID string) {
	// Get the order from the service
	coreOrder, exists := s.service.GetOrder(orderID)
	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Convert core order to API response
	items := make([]OrderItemResponse, len(coreOrder.Items))
	for i, item := range coreOrder.Items {
		items[i] = OrderItemResponse{
			Name: item.Name,
		}
	}

	apiOrder := OrderResponse{
		OrderId:      coreOrder.OrderID,
		CustomerId:   coreOrder.CustomerID,
		CreationDate: coreOrder.CreationDate,
		Status:       OrderStatus(coreOrder.Status),
		Items:        items,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiOrder)
}
