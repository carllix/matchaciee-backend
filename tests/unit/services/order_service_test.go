package services

import (
	"testing"
	"time"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/carllix/matchaciee-backend/internal/services"
	"github.com/carllix/matchaciee-backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService_CreateOrder(t *testing.T) {
	t.Run("success - create member order without customizations", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		userUUID := uuid.New()
		productUUID := uuid.New()

		user := &models.User{
			ID:       1,
			UUID:     userUUID,
			Email:    "test@example.com",
			FullName: "Test User",
			Role:     models.RoleMember,
		}

		product := &models.Product{
			ID:          1,
			UUID:        productUUID,
			Name:        "Matcha Latte",
			BasePrice:   45000,
			IsAvailable: true,
		}

		req := services.CreateOrderRequest{
			CustomerName: "Test Customer",
			Items: []services.CreateOrderItemRequest{
				{
					ProductID: productUUID,
					Quantity:  2,
				},
			},
		}

		orderNumber := "MC-260109-001"

		// Mock expectations
		mockUserRepo.On("FindByUUID", userUUID).Return(user, nil)
		mockProductRepo.On("FindByUUID", productUUID).Return(product, nil)
		mockOrderRepo.On("GenerateOrderNumber").Return(orderNumber, nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).Return(nil)
		mockOrderRepo.On("FindByUUID", mock.AnythingOfType("uuid.UUID")).Return(&models.Order{
			UUID:         uuid.New(),
			OrderNumber:  orderNumber,
			CustomerName: "Test Customer",
			Status:       models.OrderStatusPending,
			OrderSource:  models.OrderSourceMember,
			Subtotal:     90000,
			Tax:          9000,
			Total:        99000,
			Items: []models.OrderItem{
				{
					UUID:        uuid.New(),
					ProductName: "Matcha Latte",
					Quantity:    2,
					UnitPrice:   45000,
					Subtotal:    90000,
				},
			},
		}, nil)

		result, err := service.CreateOrder(userUUID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, orderNumber, result.OrderNumber)
		assert.Equal(t, "Test Customer", result.CustomerName)
		assert.Equal(t, models.OrderStatusPending, result.Status)
		assert.Equal(t, models.OrderSourceMember, result.OrderSource)
		assert.Equal(t, 90000.0, result.Subtotal)
		assert.Equal(t, 9000.0, result.Tax)
		assert.Equal(t, 99000.0, result.Total)
		assert.Len(t, result.Items, 1)

		mockUserRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("success - create member order with customizations", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		userUUID := uuid.New()
		productUUID := uuid.New()
		customizationUUID := uuid.New()

		user := &models.User{
			ID:   1,
			UUID: userUUID,
		}

		product := &models.Product{
			ID:             1,
			UUID:           productUUID,
			Name:           "Matcha Latte",
			BasePrice:      45000,
			IsAvailable:    true,
			IsCustomizable: true,
		}

		customization := &models.ProductCustomization{
			ID:                1,
			UUID:              customizationUUID,
			ProductID:         1,
			CustomizationType: "Size",
			OptionName:        "Large",
			PriceModifier:     5000,
		}

		req := services.CreateOrderRequest{
			CustomerName: "Test Customer",
			Items: []services.CreateOrderItemRequest{
				{
					ProductID: productUUID,
					Quantity:  1,
					Customizations: []services.OrderItemCustomization{
						{
							CustomizationID: customizationUUID,
							OptionName:      "Large",
						},
					},
				},
			},
		}

		orderNumber := "MC-260109-002"

		mockUserRepo.On("FindByUUID", userUUID).Return(user, nil)
		mockProductRepo.On("FindByUUID", productUUID).Return(product, nil)
		mockProductRepo.On("FindCustomizationByUUID", customizationUUID).Return(customization, nil)
		mockOrderRepo.On("GenerateOrderNumber").Return(orderNumber, nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).Return(nil)
		mockOrderRepo.On("FindByUUID", mock.AnythingOfType("uuid.UUID")).Return(&models.Order{
			UUID:         uuid.New(),
			OrderNumber:  orderNumber,
			CustomerName: "Test Customer",
			Status:       models.OrderStatusPending,
			OrderSource:  models.OrderSourceMember,
			Subtotal:     50000,
			Tax:          5000,
			Total:        55000,
			Items: []models.OrderItem{
				{
					UUID:        uuid.New(),
					ProductName: "Matcha Latte",
					Quantity:    1,
					UnitPrice:   50000, // 45000 + 5000
					Subtotal:    50000,
				},
			},
		}, nil)

		result, err := service.CreateOrder(userUUID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 50000.0, result.Subtotal)
		assert.Equal(t, 5000.0, result.Tax)
		assert.Equal(t, 55000.0, result.Total)

		mockUserRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("error - user not found", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		userUUID := uuid.New()

		req := services.CreateOrderRequest{
			CustomerName: "Test Customer",
			Items: []services.CreateOrderItemRequest{
				{
					ProductID: uuid.New(),
					Quantity:  1,
				},
			},
		}

		mockUserRepo.On("FindByUUID", userUUID).Return(nil, repositories.ErrUserNotFound)

		result, err := service.CreateOrder(userUUID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, services.ErrUserNotFound)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error - product not found", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		userUUID := uuid.New()
		productUUID := uuid.New()

		user := &models.User{
			ID:   1,
			UUID: userUUID,
		}

		req := services.CreateOrderRequest{
			CustomerName: "Test Customer",
			Items: []services.CreateOrderItemRequest{
				{
					ProductID: productUUID,
					Quantity:  1,
				},
			},
		}

		mockUserRepo.On("FindByUUID", userUUID).Return(user, nil)
		mockProductRepo.On("FindByUUID", productUUID).Return(nil, repositories.ErrProductNotFound)

		result, err := service.CreateOrder(userUUID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")

		mockUserRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - product not available", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		userUUID := uuid.New()
		productUUID := uuid.New()

		user := &models.User{
			ID:   1,
			UUID: userUUID,
		}

		product := &models.Product{
			ID:          1,
			UUID:        productUUID,
			Name:        "Matcha Latte",
			BasePrice:   45000,
			IsAvailable: false, // Not available
		}

		req := services.CreateOrderRequest{
			CustomerName: "Test Customer",
			Items: []services.CreateOrderItemRequest{
				{
					ProductID: productUUID,
					Quantity:  1,
				},
			},
		}

		mockUserRepo.On("FindByUUID", userUUID).Return(user, nil)
		mockProductRepo.On("FindByUUID", productUUID).Return(product, nil)

		result, err := service.CreateOrder(userUUID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, services.ErrProductNotAvailable)

		mockUserRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - product not customizable but customizations provided", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		userUUID := uuid.New()
		productUUID := uuid.New()

		user := &models.User{
			ID:   1,
			UUID: userUUID,
		}

		product := &models.Product{
			ID:             1,
			UUID:           productUUID,
			Name:           "Matcha Latte",
			BasePrice:      45000,
			IsAvailable:    true,
			IsCustomizable: false, // Not customizable
		}

		req := services.CreateOrderRequest{
			CustomerName: "Test Customer",
			Items: []services.CreateOrderItemRequest{
				{
					ProductID: productUUID,
					Quantity:  1,
					Customizations: []services.OrderItemCustomization{
						{
							CustomizationID: uuid.New(),
							OptionName:      "Large",
						},
					},
				},
			},
		}

		mockUserRepo.On("FindByUUID", userUUID).Return(user, nil)
		mockProductRepo.On("FindByUUID", productUUID).Return(product, nil)

		result, err := service.CreateOrder(userUUID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, services.ErrProductNotCustomizable)

		mockUserRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestOrderService_CreateGuestOrder(t *testing.T) {
	t.Run("success - create guest order", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		productUUID := uuid.New()

		product := &models.Product{
			ID:          1,
			UUID:        productUUID,
			Name:        "Matcha Latte",
			BasePrice:   45000,
			IsAvailable: true,
		}

		req := services.CreateOrderRequest{
			CustomerName: "Guest Customer",
			Items: []services.CreateOrderItemRequest{
				{
					ProductID: productUUID,
					Quantity:  1,
				},
			},
		}

		orderNumber := "MC-260109-003"

		mockProductRepo.On("FindByUUID", productUUID).Return(product, nil)
		mockOrderRepo.On("GenerateOrderNumber").Return(orderNumber, nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).Return(nil)
		mockOrderRepo.On("FindByUUID", mock.AnythingOfType("uuid.UUID")).Return(&models.Order{
			UUID:         uuid.New(),
			OrderNumber:  orderNumber,
			CustomerName: "Guest Customer",
			Status:       models.OrderStatusPending,
			OrderSource:  models.OrderSourceGuest,
			UserID:       nil, // Guest order
			Subtotal:     45000,
			Tax:          4500,
			Total:        49500,
			Items: []models.OrderItem{
				{
					UUID:        uuid.New(),
					ProductName: "Matcha Latte",
					Quantity:    1,
					UnitPrice:   45000,
					Subtotal:    45000,
				},
			},
		}, nil)

		result, err := service.CreateGuestOrder(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, orderNumber, result.OrderNumber)
		assert.Equal(t, models.OrderSourceGuest, result.OrderSource)
		assert.Nil(t, result.User) // Guest order has no user

		mockProductRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetByUUID(t *testing.T) {
	t.Run("success - get order by UUID", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		orderUUID := uuid.New()

		order := &models.Order{
			UUID:         orderUUID,
			OrderNumber:  "MC-260109-001",
			CustomerName: "Test Customer",
			Status:       models.OrderStatusPending,
			OrderSource:  models.OrderSourceMember,
			Subtotal:     45000,
			Tax:          4500,
			Total:        49500,
			Items: []models.OrderItem{
				{
					UUID:        uuid.New(),
					ProductName: "Matcha Latte",
					Quantity:    1,
					UnitPrice:   45000,
					Subtotal:    45000,
				},
			},
		}

		mockOrderRepo.On("FindByUUID", orderUUID).Return(order, nil)

		result, err := service.GetByUUID(orderUUID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, orderUUID, result.ID)
		assert.Equal(t, "MC-260109-001", result.OrderNumber)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("error - order not found", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		orderUUID := uuid.New()

		mockOrderRepo.On("FindByUUID", orderUUID).Return(nil, repositories.ErrOrderNotFound)

		result, err := service.GetByUUID(orderUUID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, services.ErrOrderNotFound)

		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetMyOrders(t *testing.T) {
	t.Run("success - get user orders", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		userUUID := uuid.New()

		user := &models.User{
			ID:   1,
			UUID: userUUID,
		}

		orders := []models.Order{
			{
				UUID:         uuid.New(),
				OrderNumber:  "MC-260109-001",
				CustomerName: "Test Customer",
				Status:       models.OrderStatusPending,
				Subtotal:     45000,
				Tax:          4500,
				Total:        49500,
				Items:        []models.OrderItem{},
			},
		}

		mockUserRepo.On("FindByUUID", userUUID).Return(user, nil)
		mockOrderRepo.On("FindByUserID", uint(1), 10, 0).Return(orders, int64(1), nil)

		result, err := service.GetMyOrders(userUUID, 1, 10)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Orders, 1)
		assert.Equal(t, int64(1), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.Limit)

		mockUserRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_UpdateOrderStatus(t *testing.T) {
	t.Run("success - valid status transition from pending to preparing", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		orderUUID := uuid.New()

		order := &models.Order{
			ID:          1,
			UUID:        orderUUID,
			OrderNumber: "MC-260109-001",
			Status:      models.OrderStatusPending,
		}

		updatedOrder := &models.Order{
			ID:          1,
			UUID:        orderUUID,
			OrderNumber: "MC-260109-001",
			Status:      models.OrderStatusPreparing,
			Items:       []models.OrderItem{},
		}

		mockOrderRepo.On("FindByUUID", orderUUID).Return(order, nil).Once()
		mockOrderRepo.On("UpdateStatus", uint(1), models.OrderStatusPreparing).Return(nil)
		mockOrderRepo.On("FindByUUID", orderUUID).Return(updatedOrder, nil).Once()

		result, err := service.UpdateOrderStatus(orderUUID, models.OrderStatusPreparing)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.OrderStatusPreparing, result.Status)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("success - valid status transition from preparing to ready", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		orderUUID := uuid.New()

		order := &models.Order{
			ID:          1,
			UUID:        orderUUID,
			OrderNumber: "MC-260109-001",
			Status:      models.OrderStatusPreparing,
		}

		updatedOrder := &models.Order{
			ID:          1,
			UUID:        orderUUID,
			OrderNumber: "MC-260109-001",
			Status:      models.OrderStatusReady,
			Items:       []models.OrderItem{},
		}

		mockOrderRepo.On("FindByUUID", orderUUID).Return(order, nil).Once()
		mockOrderRepo.On("UpdateStatus", uint(1), models.OrderStatusReady).Return(nil)
		mockOrderRepo.On("FindByUUID", orderUUID).Return(updatedOrder, nil).Once()

		result, err := service.UpdateOrderStatus(orderUUID, models.OrderStatusReady)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.OrderStatusReady, result.Status)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("success - valid status transition from ready to completed", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		orderUUID := uuid.New()

		order := &models.Order{
			ID:          1,
			UUID:        orderUUID,
			OrderNumber: "MC-260109-001",
			Status:      models.OrderStatusReady,
		}

		completedTime := time.Now()
		updatedOrder := &models.Order{
			ID:          1,
			UUID:        orderUUID,
			OrderNumber: "MC-260109-001",
			Status:      models.OrderStatusCompleted,
			CompletedAt: &completedTime,
			Items:       []models.OrderItem{},
		}

		mockOrderRepo.On("FindByUUID", orderUUID).Return(order, nil).Once()
		mockOrderRepo.On("UpdateStatus", uint(1), models.OrderStatusCompleted).Return(nil)
		mockOrderRepo.On("FindByUUID", orderUUID).Return(updatedOrder, nil).Once()

		result, err := service.UpdateOrderStatus(orderUUID, models.OrderStatusCompleted)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.OrderStatusCompleted, result.Status)
		assert.NotNil(t, result.CompletedAt)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("error - invalid status transition", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		orderUUID := uuid.New()

		order := &models.Order{
			ID:          1,
			UUID:        orderUUID,
			OrderNumber: "MC-260109-001",
			Status:      models.OrderStatusPending,
		}

		mockOrderRepo.On("FindByUUID", orderUUID).Return(order, nil)

		// Try to transition from pending to ready (invalid - should go to preparing first)
		result, err := service.UpdateOrderStatus(orderUUID, models.OrderStatusReady)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, services.ErrInvalidStatusTransition)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("error - transition from completed status", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		orderUUID := uuid.New()

		order := &models.Order{
			ID:          1,
			UUID:        orderUUID,
			OrderNumber: "MC-260109-001",
			Status:      models.OrderStatusCompleted,
		}

		mockOrderRepo.On("FindByUUID", orderUUID).Return(order, nil)

		// Try to transition from completed to any other status (invalid)
		result, err := service.UpdateOrderStatus(orderUUID, models.OrderStatusReady)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, services.ErrInvalidStatusTransition)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("error - order not found", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		orderUUID := uuid.New()

		mockOrderRepo.On("FindByUUID", orderUUID).Return(nil, repositories.ErrOrderNotFound)

		result, err := service.UpdateOrderStatus(orderUUID, models.OrderStatusPreparing)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, services.ErrOrderNotFound)

		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetAllOrders(t *testing.T) {
	t.Run("success - get all orders with filters", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		status := models.OrderStatusPending
		filters := repositories.OrderFilters{
			Status: &status,
		}

		orders := []models.Order{
			{
				UUID:         uuid.New(),
				OrderNumber:  "MC-260109-001",
				CustomerName: "Customer 1",
				Status:       models.OrderStatusPending,
				Items:        []models.OrderItem{},
			},
			{
				UUID:         uuid.New(),
				OrderNumber:  "MC-260109-002",
				CustomerName: "Customer 2",
				Status:       models.OrderStatusPending,
				Items:        []models.OrderItem{},
			},
		}

		mockOrderRepo.On("FindAll", filters, 20, 0).Return(orders, int64(2), nil)

		result, err := service.GetAllOrders(filters, 1, 20)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Orders, 2)
		assert.Equal(t, int64(2), result.Total)

		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_CalculateTotals(t *testing.T) {
	t.Run("correct calculation - single item without customizations", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		productUUID := uuid.New()

		product := &models.Product{
			ID:          1,
			UUID:        productUUID,
			Name:        "Matcha Latte",
			BasePrice:   45000,
			IsAvailable: true,
		}

		user := &models.User{
			ID:   1,
			UUID: uuid.New(),
		}

		req := services.CreateOrderRequest{
			CustomerName: "Test",
			Items: []services.CreateOrderItemRequest{
				{
					ProductID: productUUID,
					Quantity:  2,
				},
			},
		}

		mockUserRepo.On("FindByUUID", mock.Anything).Return(user, nil)
		mockProductRepo.On("FindByUUID", productUUID).Return(product, nil)
		mockOrderRepo.On("GenerateOrderNumber").Return("MC-260109-001", nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).
			Run(func(args mock.Arguments) {
				order := args.Get(0).(*models.Order)
				// Subtotal = 45000 * 2 = 90000
				// Tax = 90000 * 0.1 = 9000
				// Total = 90000 + 9000 = 99000
				assert.Equal(t, 90000.0, order.Subtotal)
				assert.Equal(t, 9000.0, order.Tax)
				assert.Equal(t, 99000.0, order.Total)
			}).
			Return(nil)
		mockOrderRepo.On("FindByUUID", mock.AnythingOfType("uuid.UUID")).Return(&models.Order{
			UUID:     uuid.New(),
			Subtotal: 90000,
			Tax:      9000,
			Total:    99000,
			Items:    []models.OrderItem{},
		}, nil)

		_, err := service.CreateOrder(user.UUID, req)

		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("correct calculation - with customizations", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		service := services.NewOrderService(mockOrderRepo, mockProductRepo, mockUserRepo)

		productUUID := uuid.New()
		customizationUUID := uuid.New()

		product := &models.Product{
			ID:             1,
			UUID:           productUUID,
			Name:           "Matcha Latte",
			BasePrice:      45000,
			IsAvailable:    true,
			IsCustomizable: true,
		}

		customization := &models.ProductCustomization{
			ID:            1,
			UUID:          customizationUUID,
			ProductID:     1,
			PriceModifier: 5000,
		}

		user := &models.User{
			ID:   1,
			UUID: uuid.New(),
		}

		req := services.CreateOrderRequest{
			CustomerName: "Test",
			Items: []services.CreateOrderItemRequest{
				{
					ProductID: productUUID,
					Quantity:  1,
					Customizations: []services.OrderItemCustomization{
						{CustomizationID: customizationUUID},
					},
				},
			},
		}

		mockUserRepo.On("FindByUUID", mock.Anything).Return(user, nil)
		mockProductRepo.On("FindByUUID", productUUID).Return(product, nil)
		mockProductRepo.On("FindCustomizationByUUID", customizationUUID).Return(customization, nil)
		mockOrderRepo.On("GenerateOrderNumber").Return("MC-260109-001", nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).
			Run(func(args mock.Arguments) {
				order := args.Get(0).(*models.Order)
				// Subtotal = (45000 + 5000) * 1 = 50000
				// Tax = 50000 * 0.1 = 5000
				// Total = 50000 + 5000 = 55000
				assert.Equal(t, 50000.0, order.Subtotal)
				assert.Equal(t, 5000.0, order.Tax)
				assert.Equal(t, 55000.0, order.Total)
			}).
			Return(nil)
		mockOrderRepo.On("FindByUUID", mock.AnythingOfType("uuid.UUID")).Return(&models.Order{
			UUID:     uuid.New(),
			Subtotal: 50000,
			Tax:      5000,
			Total:    55000,
			Items:    []models.OrderItem{},
		}, nil)

		_, err := service.CreateOrder(user.UUID, req)

		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
	})
}
