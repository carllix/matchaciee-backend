package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

var (
	ErrOrderNotFound           = errors.New("order not found")
	ErrUserNotFound            = errors.New("user not found")
	ErrProductNotAvailable     = errors.New("product not available")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrProductNotCustomizable  = errors.New("product is not customizable")
	ErrInvalidCustomization    = errors.New("customization does not belong to product")
)

type CreateOrderRequest struct {
	CustomerName string                   `json:"customer_name" validate:"required,min=2,max=255"`
	Notes        *string                  `json:"notes,omitempty" validate:"omitempty,max=500"`
	Items        []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID      uuid.UUID                `json:"product_id" validate:"required"`
	Quantity       int                      `json:"quantity" validate:"required,min=1,max=100"`
	Notes          *string                  `json:"notes,omitempty" validate:"omitempty,max=200"`
	Customizations []OrderItemCustomization `json:"customizations,omitempty" validate:"omitempty,dive"`
}

type OrderItemCustomization struct {
	CustomizationID uuid.UUID `json:"customization_id" validate:"required"`
	OptionName      string    `json:"option_name" validate:"required"`
}

type UpdateOrderStatusRequest struct {
	Status models.OrderStatus `json:"status" validate:"required,oneof=pending preparing ready completed cancelled"`
}

type OrderResponse struct {
	ID           uuid.UUID           `json:"id"`
	OrderNumber  string              `json:"order_number"`
	CustomerName string              `json:"customer_name"`
	Status       models.OrderStatus  `json:"status"`
	OrderSource  models.OrderSource  `json:"order_source"`
	Subtotal     float64             `json:"subtotal"`
	Tax          float64             `json:"tax"`
	Total        float64             `json:"total"`
	Notes        *string             `json:"notes,omitempty"`
	Items        []OrderItemResponse `json:"items"`
	User         *UserSummary        `json:"user,omitempty"`
	CreatedAt    string              `json:"created_at"`
	CompletedAt  *string             `json:"completed_at,omitempty"`
}

type OrderItemResponse struct {
	ID             uuid.UUID      `json:"id"`
	ProductName    string         `json:"product_name"`
	Quantity       int            `json:"quantity"`
	UnitPrice      float64        `json:"unit_price"`
	Subtotal       float64        `json:"subtotal"`
	Notes          *string        `json:"notes,omitempty"`
	Customizations datatypes.JSON `json:"customizations,omitempty"`
}

type UserSummary struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

type OrderService interface {
	CreateOrder(userUUID uuid.UUID, req CreateOrderRequest) (*OrderResponse, error)
	CreateGuestOrder(req CreateOrderRequest) (*OrderResponse, error)
	GetByUUID(orderUUID uuid.UUID) (*OrderResponse, error)
	GetByOrderNumber(orderNumber string) (*OrderResponse, error)
	GetMyOrders(userUUID uuid.UUID, page, limit int) (*OrderListResponse, error)
	GetAllOrders(filters repositories.OrderFilters, page, limit int) (*OrderListResponse, error)
	UpdateOrderStatus(orderUUID uuid.UUID, status models.OrderStatus) (*OrderResponse, error)
}

type orderService struct {
	orderRepo   repositories.OrderRepository
	productRepo repositories.ProductRepository
	userRepo    repositories.UserRepository
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	userRepo repositories.UserRepository,
) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

func (s *orderService) CreateOrder(userUUID uuid.UUID, req CreateOrderRequest) (*OrderResponse, error) {
	// Get user by UUID to get internal ID
	user, err := s.userRepo.FindByUUID(userUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Validate and fetch all products
	products, customizationsMap, err := s.validateAndFetchProducts(req.Items)
	if err != nil {
		return nil, err
	}

	// Calculate totals and build order items
	subtotal, orderItems := s.calculateOrderTotals(req.Items, products, customizationsMap)
	tax := subtotal * 0.10 // 10% tax
	total := subtotal + tax

	// Generate order number
	orderNumber, err := s.orderRepo.GenerateOrderNumber()
	if err != nil {
		return nil, err
	}

	// Build order object
	order := &models.Order{
		OrderNumber:  orderNumber,
		UserID:       &user.ID,
		CustomerName: req.CustomerName,
		Notes:        req.Notes,
		Status:       models.OrderStatusPending,
		OrderSource:  models.OrderSourceMember,
		Subtotal:     subtotal,
		Tax:          tax,
		Total:        total,
	}

	// Create order
	err = s.orderRepo.Create(order, orderItems)
	if err != nil {
		return nil, err
	}

	// Fetch complete order with relations
	createdOrder, err := s.orderRepo.FindByUUID(order.UUID)
	if err != nil {
		return nil, err
	}

	return s.toOrderResponse(createdOrder, false), nil
}

func (s *orderService) CreateGuestOrder(req CreateOrderRequest) (*OrderResponse, error) {
	// Validate and fetch all products
	products, customizationsMap, err := s.validateAndFetchProducts(req.Items)
	if err != nil {
		return nil, err
	}

	// Calculate totals and build order items
	subtotal, orderItems := s.calculateOrderTotals(req.Items, products, customizationsMap)
	tax := subtotal * 0.10 // 10% tax
	total := subtotal + tax

	// Generate order number
	orderNumber, err := s.orderRepo.GenerateOrderNumber()
	if err != nil {
		return nil, err
	}

	// Build order object
	order := &models.Order{
		OrderNumber:  orderNumber,
		UserID:       nil,
		CustomerName: req.CustomerName,
		Notes:        req.Notes,
		Status:       models.OrderStatusPending,
		OrderSource:  models.OrderSourceGuest,
		Subtotal:     subtotal,
		Tax:          tax,
		Total:        total,
	}

	// Create order
	err = s.orderRepo.Create(order, orderItems)
	if err != nil {
		return nil, err
	}

	// Fetch complete order with relations
	createdOrder, err := s.orderRepo.FindByUUID(order.UUID)
	if err != nil {
		return nil, err
	}

	return s.toOrderResponse(createdOrder, false), nil
}

func (s *orderService) GetByUUID(orderUUID uuid.UUID) (*OrderResponse, error) {
	order, err := s.orderRepo.FindByUUID(orderUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return s.toOrderResponse(order, false), nil
}

func (s *orderService) GetByOrderNumber(orderNumber string) (*OrderResponse, error) {
	order, err := s.orderRepo.FindByOrderNumber(orderNumber)
	if err != nil {
		if errors.Is(err, repositories.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Include user details (admin access)
	return s.toOrderResponse(order, true), nil
}

func (s *orderService) GetMyOrders(userUUID uuid.UUID, page, limit int) (*OrderListResponse, error) {
	// Get user by UUID to get internal ID
	user, err := s.userRepo.FindByUUID(userUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get orders
	orders, total, err := s.orderRepo.FindByUserID(user.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert to response
	orderResponses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = *s.toOrderResponse(&order, false)
	}

	return &OrderListResponse{
		Orders: orderResponses,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

func (s *orderService) GetAllOrders(filters repositories.OrderFilters, page, limit int) (*OrderListResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get orders
	orders, total, err := s.orderRepo.FindAll(filters, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert to response
	orderResponses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = *s.toOrderResponse(&order, true)
	}

	return &OrderListResponse{
		Orders: orderResponses,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

func (s *orderService) UpdateOrderStatus(orderUUID uuid.UUID, status models.OrderStatus) (*OrderResponse, error) {
	// Fetch order
	order, err := s.orderRepo.FindByUUID(orderUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Validate status transition
	if !s.isValidStatusTransition(order.Status, status) {
		return nil, ErrInvalidStatusTransition
	}

	// Update status
	err = s.orderRepo.UpdateStatus(order.ID, status)
	if err != nil {
		return nil, err
	}

	// Fetch updated order
	updatedOrder, err := s.orderRepo.FindByUUID(orderUUID)
	if err != nil {
		return nil, err
	}

	return s.toOrderResponse(updatedOrder, true), nil
}

func (s *orderService) validateAndFetchProducts(items []CreateOrderItemRequest) (
	map[uuid.UUID]*models.Product,
	map[uuid.UUID]map[uuid.UUID]*models.ProductCustomization,
	error,
) {
	products := make(map[uuid.UUID]*models.Product)
	customizationsMap := make(map[uuid.UUID]map[uuid.UUID]*models.ProductCustomization)

	for _, item := range items {
		// Fetch product
		product, err := s.productRepo.FindByUUID(item.ProductID)
		if err != nil {
			if errors.Is(err, repositories.ErrProductNotFound) {
				return nil, nil, fmt.Errorf("product %s not found", item.ProductID)
			}
			return nil, nil, err
		}

		// Validate product is available
		if !product.IsAvailable {
			return nil, nil, fmt.Errorf("%w: %s", ErrProductNotAvailable, product.Name)
		}

		products[item.ProductID] = product

		// Validate customizations
		if len(item.Customizations) > 0 {
			if !product.IsCustomizable {
				return nil, nil, fmt.Errorf("%w: %s", ErrProductNotCustomizable, product.Name)
			}

			customMap := make(map[uuid.UUID]*models.ProductCustomization)
			for _, customReq := range item.Customizations {
				custom, err := s.productRepo.FindCustomizationByUUID(customReq.CustomizationID)
				if err != nil {
					return nil, nil, fmt.Errorf("customization not found: %s", customReq.CustomizationID)
				}

				// Verify customization belongs to product
				if custom.ProductID != product.ID {
					return nil, nil, fmt.Errorf("%w: %s", ErrInvalidCustomization, customReq.OptionName)
				}

				customMap[customReq.CustomizationID] = custom
			}
			customizationsMap[item.ProductID] = customMap
		}
	}

	return products, customizationsMap, nil
}

func (s *orderService) calculateOrderTotals(
	items []CreateOrderItemRequest,
	products map[uuid.UUID]*models.Product,
	customizationsMap map[uuid.UUID]map[uuid.UUID]*models.ProductCustomization,
) (float64, []models.OrderItem) {
	var subtotal float64
	var orderItems []models.OrderItem

	for _, item := range items {
		product := products[item.ProductID]

		// Calculate unit price (base + customizations)
		unitPrice := product.BasePrice
		var customizationsJSON datatypes.JSON

		if customMap, exists := customizationsMap[item.ProductID]; exists && len(customMap) > 0 {
			customizations := make([]map[string]any, 0)
			for _, custom := range customMap {
				unitPrice += custom.PriceModifier
				customizations = append(customizations, map[string]any{
					"customization_type": custom.CustomizationType,
					"option_name":        custom.OptionName,
					"price_modifier":     custom.PriceModifier,
				})
			}
			// Marshal customizations to JSON
			var err error
			customizationsJSON, err = json.Marshal(customizations)
			if err != nil {
				customizationsJSON = datatypes.JSON("{}")
			}
		}

		itemSubtotal := unitPrice * float64(item.Quantity)
		subtotal += itemSubtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:      &product.ID,
			ProductName:    product.Name,
			Quantity:       item.Quantity,
			UnitPrice:      unitPrice,
			Subtotal:       itemSubtotal,
			Notes:          item.Notes,
			Customizations: customizationsJSON,
		})
	}

	return subtotal, orderItems
}

func (s *orderService) isValidStatusTransition(current, new models.OrderStatus) bool {
	validTransitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusPending:   {models.OrderStatusPreparing, models.OrderStatusCancelled},
		models.OrderStatusPreparing: {models.OrderStatusReady},
		models.OrderStatusReady:     {models.OrderStatusCompleted},
	}

	allowed, exists := validTransitions[current]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == new {
			return true
		}
	}
	return false
}

func (s *orderService) toOrderResponse(order *models.Order, includeUser bool) *OrderResponse {
	itemResponses := make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		itemResponses[i] = OrderItemResponse{
			ID:             item.UUID,
			ProductName:    item.ProductName,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			Subtotal:       item.Subtotal,
			Notes:          item.Notes,
			Customizations: item.Customizations,
		}
	}

	var userSummary *UserSummary
	if includeUser && order.User != nil {
		userSummary = &UserSummary{
			ID:       order.User.UUID,
			FullName: order.User.FullName,
			Email:    order.User.Email,
		}
	}

	var completedAt *string
	if order.CompletedAt != nil {
		completedAtStr := order.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		completedAt = &completedAtStr
	}

	return &OrderResponse{
		ID:           order.UUID,
		OrderNumber:  order.OrderNumber,
		CustomerName: order.CustomerName,
		Status:       order.Status,
		OrderSource:  order.OrderSource,
		Subtotal:     order.Subtotal,
		Tax:          order.Tax,
		Total:        order.Total,
		Notes:        order.Notes,
		Items:        itemResponses,
		User:         userSummary,
		CreatedAt:    order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CompletedAt:  completedAt,
	}
}
