package handlers

import (
	"errors"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/carllix/matchaciee-backend/internal/services"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// Get user UUID from context
	userUUIDValue := c.Locals("userUUID")
	if userUUIDValue == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userUUID, ok := userUUIDValue.(uuid.UUID)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user UUID")
	}

	// Parse request
	var req services.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Create order
	order, err := h.orderService.CreateOrder(userUUID, req)
	if err != nil {
		if errors.Is(err, services.ErrProductNotAvailable) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
		if errors.Is(err, services.ErrProductNotCustomizable) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
		if errors.Is(err, services.ErrInvalidCustomization) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
		if errors.Is(err, services.ErrUserNotFound) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "User not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create order")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, order)
}

// POST /api/v1/orders/guest
func (h *OrderHandler) CreateGuestOrder(c *fiber.Ctx) error {
	var req services.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	order, err := h.orderService.CreateGuestOrder(req)
	if err != nil {
		if errors.Is(err, services.ErrProductNotAvailable) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
		if errors.Is(err, services.ErrProductNotCustomizable) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
		if errors.Is(err, services.ErrInvalidCustomization) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create guest order")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, order)
}

// GET /api/v1/orders/track/:uuid
func (h *OrderHandler) TrackGuestOrder(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")

	orderUUID, err := uuid.Parse(uuidParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid order ID format")
	}

	order, err := h.orderService.GetByUUID(orderUUID)
	if err != nil {
		if errors.Is(err, services.ErrOrderNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get order")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, order)
}

// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	idParam := c.Params("id")

	orderUUID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid order ID format")
	}

	// Get user info from context
	userUUIDValue := c.Locals("userUUID")
	roleValue := c.Locals("role")

	order, err := h.orderService.GetByUUID(orderUUID)
	if err != nil {
		if errors.Is(err, services.ErrOrderNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get order")
	}

	// Authorization check: member can only view own orders
	role, ok := roleValue.(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid role")
	}
	if role == string(models.RoleMember) {
		userUUID, ok := userUUIDValue.(uuid.UUID)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user UUID")
		}
		if order.User == nil || order.User.ID != userUUID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied")
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, order)
}

// GET /api/v1/orders/number/:number
func (h *OrderHandler) GetOrderByNumber(c *fiber.Ctx) error {
	orderNumber := c.Params("number")

	order, err := h.orderService.GetByOrderNumber(orderNumber)
	if err != nil {
		if errors.Is(err, services.ErrOrderNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get order")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, order)
}

// GET /api/v1/orders/me
func (h *OrderHandler) GetMyOrders(c *fiber.Ctx) error {
	userUUIDValue := c.Locals("userUUID")
	if userUUIDValue == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userUUID, ok := userUUIDValue.(uuid.UUID)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user UUID")
	}

	// Pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orders, err := h.orderService.GetMyOrders(userUUID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get orders")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, orders)
}

// GET /api/v1/orders
func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	// Pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Filters
	var filters repositories.OrderFilters

	if statusParam := c.Query("status"); statusParam != "" {
		status := models.OrderStatus(statusParam)
		filters.Status = &status
	}

	if sourceParam := c.Query("source"); sourceParam != "" {
		source := models.OrderSource(sourceParam)
		filters.OrderSource = &source
	}

	orders, err := h.orderService.GetAllOrders(filters, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get orders")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, orders)
}

// PUT /api/v1/orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")

	orderUUID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid order ID format")
	}

	var req services.UpdateOrderStatusRequest
	if err = c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	order, err := h.orderService.UpdateOrderStatus(orderUUID, req.Status)
	if err != nil {
		if errors.Is(err, services.ErrOrderNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found")
		}
		if errors.Is(err, services.ErrInvalidStatusTransition) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid status transition")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update order status")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, order)
}
