package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	domain "reference-service-go/internal/domain/order"
	ports "reference-service-go/internal/ports/order"
)

// API is a thin HTTP adapter that talks to the use case service
type API struct {
	service ports.Service
}

func NewAPI(service ports.Service) *API { return &API{service: service} }

// GetOrders returns a paginated list of orders with optional filtering
func (api *API) GetOrders(w http.ResponseWriter, _ *http.Request, params GetOrdersParams) {
	limit := 50
	offset := 0
	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}

	var customerID *string
	if params.CustomerId != nil {
		id := params.CustomerId.String()
		customerID = &id
	}

	orders := api.service.GetOrders(customerID, limit, offset)

	page := OrdersPage{
		Data: make([]OrderResponse, 0, len(orders)),
	}
	// Fill pagination metadata per OpenAPI
	if v := len(orders); v >= 0 { // always true, but explicit
		vv := v
		page.Total = &vv
	}
	if params.Limit != nil {
		page.Limit = params.Limit
	} else {
		v := limit
		page.Limit = &v
	}
	if params.Offset != nil {
		page.Offset = params.Offset
	} else {
		v := offset
		page.Offset = &v
	}

	for _, order := range orders {
		items := make([]OrderItemResponse, len(order.Items))
		for i, item := range order.Items {
			items[i] = OrderItemResponse{Name: item.Name}
		}
		page.Data = append(page.Data, OrderResponse{
			OrderId:    order.OrderID,
			CustomerId: openapi_types.UUID(uuid.MustParse(order.CustomerID)),
			CreatedAt:  order.CreationDate,
			Status:     OrderStatus(order.Status),
			Items:      items,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(page)
}

// PostOrders creates a new order
func (api *API) PostOrders(w http.ResponseWriter, r *http.Request) {
	var request OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	items := make([]domain.OrderItem, len(request.Items))
	for i, item := range request.Items {
		items[i] = domain.OrderItem{Name: item.Name}
	}
	domainRequest := domain.OrderRequest{
		CustomerID: request.CustomerId.String(),
		Items:      items,
		Country:    request.Country,
	}

	createdOrder, err := api.service.CreateOrder(domainRequest)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCountry) {
			code := http.StatusBadRequest
			detail := "country must be exactly two uppercase letters [A-Z]{2} and already uppercase"
			problem := Problem{
				Type:   "https://example.com/problems/invalid-country",
				Title:  "Invalid country",
				Status: &code,
				Detail: &detail,
			}
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(code)
			_ = json.NewEncoder(w).Encode(problem)
			return
		}
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	responseItems := make([]OrderItemResponse, len(createdOrder.Items))
	for i, item := range createdOrder.Items {
		responseItems[i] = OrderItemResponse{Name: item.Name}
	}
	response := OrderResponse{
		OrderId:    createdOrder.OrderID,
		CustomerId: openapi_types.UUID(uuid.MustParse(createdOrder.CustomerID)),
		CreatedAt:  createdOrder.CreationDate,
		Status:     OrderStatus(createdOrder.Status),
		Items:      responseItems,
	}

	// Set Location header per OpenAPI
	w.Header().Set("Location", fmt.Sprintf("/v1/orders/%s", createdOrder.OrderID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

// GetOrder returns one order
func (api *API) GetOrder(w http.ResponseWriter, _ *http.Request, orderID string) {
	order, found := api.service.GetOrder(orderID)
	if !found {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	items := make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemResponse{Name: item.Name}
	}
	response := OrderResponse{
		OrderId:    order.OrderID,
		CustomerId: uuid.MustParse(order.CustomerID),
		CreatedAt:  order.CreationDate,
		Status:     order.Status,
		Items:      items,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
