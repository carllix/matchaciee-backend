package handlers

import (
	"errors"

	_ "github.com/carllix/matchaciee-backend/docs"
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

// CreateOrder godoc
// @Summary Create an order (authenticated)
// @Description Create a new order for an authenticated member user
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body docs.CreateOrderRequest true "Order details"
// @Success 201 {object} docs.OrderSuccessResponse "Order created successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error, product not available, or invalid customization"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /orders [post]
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

// CreateGuestOrder godoc
// @Summary Create a guest order
// @Description Create a new order without authentication. Order can be tracked via the returned order UUID.
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body docs.CreateOrderRequest true "Order details"
// @Success 201 {object} docs.OrderSuccessResponse "Order created successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error, product not available, or invalid customization"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /orders/guest [post]
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

// TrackGuestOrder godoc
// @Summary Track a guest order
// @Description Track an order by its UUID. This is public and used for guest order tracking.
// @Tags Orders
// @Accept json
// @Produce json
// @Param uuid path string true "Order UUID"
// @Success 200 {object} docs.OrderSuccessResponse "Order retrieved successfully"
// @Failure 400 {object} docs.SwaggerErrorResponse "Invalid order ID format"
// @Failure 404 {object} docs.SwaggerErrorResponse "Order not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /orders/track/{uuid} [get]
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

// GetOrder godoc
// @Summary Get order by ID
// @Description Get a single order by its UUID. Members can only view their own orders, Admin/Barista can view any order.
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order UUID"
// @Success 200 {object} docs.OrderSuccessResponse "Order retrieved successfully"
// @Failure 400 {object} docs.SwaggerErrorResponse "Invalid order ID format"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Access denied"
// @Failure 404 {object} docs.SwaggerErrorResponse "Order not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /orders/{id} [get]
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

// GetOrderByNumber godoc
// @Summary Get order by order number
// @Description Get a single order by its order number (e.g., MC-250107-001). Admin only.
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param number path string true "Order number (e.g., MC-250107-001)"
// @Success 200 {object} docs.OrderSuccessResponse "Order retrieved successfully"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} docs.SwaggerErrorResponse "Order not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /orders/number/{number} [get]
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

// GetMyOrders godoc
// @Summary Get my orders
// @Description Get a paginated list of the authenticated user's orders
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query integer false "Page number" default(1)
// @Param limit query integer false "Items per page (max 100)" default(10)
// @Success 200 {object} docs.OrdersSuccessResponse "Orders retrieved successfully"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /orders/me [get]
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

// GetAllOrders godoc
// @Summary Get all orders
// @Description Get a paginated list of all orders with optional filtering. Admin/Barista only.
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query integer false "Page number" default(1)
// @Param limit query integer false "Items per page (max 100)" default(20)
// @Param status query string false "Filter by order status" Enums(pending, preparing, ready, completed, cancelled)
// @Param source query string false "Filter by order source" Enums(guest, member, kiosk)
// @Success 200 {object} docs.OrdersSuccessResponse "Orders retrieved successfully"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin/Barista only"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /orders [get]
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

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an order. Admin/Barista only. Status transitions must follow: pending -> preparing -> ready -> completed, or pending -> cancelled.
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order UUID"
// @Param request body docs.UpdateOrderStatusRequest true "New order status"
// @Success 200 {object} docs.OrderSuccessResponse "Order status updated successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error, invalid ID format, or invalid status transition"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin/Barista only"
// @Failure 404 {object} docs.SwaggerErrorResponse "Order not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /orders/{id}/status [put]
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
