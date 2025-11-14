package order

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	dom "reference-service-go/internal/domain/order"
	uc "reference-service-go/internal/usecase/order"
)

// API is a thin HTTP adapter that talks to the use case service
type API struct {
	service *uc.Service
}

func NewAPI(service *uc.Service) *API { return &API{service: service} }

// GetOrders returns a list of orders with optional filtering
func (h *API) GetOrders(w http.ResponseWriter, r *http.Request, params GetOrdersParams) {
	limit := 50
	offset := 0
	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}

	var customerIDStr *string
	if params.CustomerId != nil {
		id := params.CustomerId.String()
		customerIDStr = &id
	}

	orders := h.service.GetOrders(customerIDStr, limit, offset)

	var result OrdersResponse
	for _, o := range orders {
		items := make([]OrderItemResponse, len(o.Items))
		for i, it := range o.Items {
			items[i] = OrderItemResponse{Name: it.Name}
		}
		result = append(result, OrderResponse{
			OrderId:      o.OrderID,
			CustomerId:   openapi_types.UUID(uuid.MustParse(o.CustomerID)),
			CreationDate: o.CreationDate,
			Status:       OrderStatus(o.Status),
			Items:        items,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

// PostOrders creates a new order
func (h *API) PostOrders(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	items := make([]dom.OrderItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = dom.OrderItem{Name: it.Name}
	}
	coreReq := dom.OrderRequest{
		CustomerID: req.CustomerId.String(),
		Items:      items,
	}

	ord, err := h.service.CreateOrder(coreReq)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	respItems := make([]OrderItemResponse, len(ord.Items))
	for i, it := range ord.Items {
		respItems[i] = OrderItemResponse{Name: it.Name}
	}
	resp := OrderResponse{
		OrderId:      ord.OrderID,
		CustomerId:   openapi_types.UUID(uuid.MustParse(ord.CustomerID)),
		CreationDate: ord.CreationDate,
		Status:       OrderStatus(ord.Status),
		Items:        respItems,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// GetOrder returns one order
func (h *API) GetOrder(w http.ResponseWriter, r *http.Request, orderID string) {
	ord, ok := h.service.GetOrder(orderID)
	if !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	items := make([]OrderItemResponse, len(ord.Items))
	for i, it := range ord.Items {
		items[i] = OrderItemResponse{Name: it.Name}
	}
	resp := OrderResponse{
		OrderId:      ord.OrderID,
		CustomerId:   openapi_types.UUID(uuid.MustParse(ord.CustomerID)),
		CreationDate: ord.CreationDate,
		Status:       OrderStatus(ord.Status),
		Items:        items,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
