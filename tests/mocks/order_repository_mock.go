package mocks

import (
	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *models.Order, items []models.OrderItem) error {
	args := m.Called(order, items)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByID(id uint) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByUUID(orderUUID uuid.UUID) (*models.Order, error) {
	args := m.Called(orderUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByOrderNumber(orderNumber string) (*models.Order, error) {
	args := m.Called(orderNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByUserID(userID uint, limit, offset int) ([]models.Order, int64, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) FindAll(filters repositories.OrderFilters, limit, offset int) ([]models.Order, int64, error) {
	args := m.Called(filters, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) UpdateStatus(orderID uint, status models.OrderStatus) error {
	args := m.Called(orderID, status)
	return args.Error(0)
}

func (m *MockOrderRepository) GenerateOrderNumber() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
