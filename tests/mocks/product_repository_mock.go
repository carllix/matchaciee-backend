package mocks

import (
	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) FindByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	product, ok := args.Get(0).(*models.Product)
	if !ok {
		return nil, args.Error(1)
	}
	return product, args.Error(1)
}

func (m *MockProductRepository) FindByUUID(uuid uuid.UUID) (*models.Product, error) {
	args := m.Called(uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	product, ok := args.Get(0).(*models.Product)
	if !ok {
		return nil, args.Error(1)
	}
	return product, args.Error(1)
}

func (m *MockProductRepository) FindByUUIDIncludingDeleted(uuid uuid.UUID) (*models.Product, error) {
	args := m.Called(uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	product, ok := args.Get(0).(*models.Product)
	if !ok {
		return nil, args.Error(1)
	}
	return product, args.Error(1)
}

func (m *MockProductRepository) FindBySlug(slug string) (*models.Product, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	product, ok := args.Get(0).(*models.Product)
	if !ok {
		return nil, args.Error(1)
	}
	return product, args.Error(1)
}

func (m *MockProductRepository) FindAll(includeDeleted bool, isAvailable *bool, categoryID *uint) ([]models.Product, error) {
	args := m.Called(includeDeleted, isAvailable, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	products, ok := args.Get(0).([]models.Product)
	if !ok {
		return nil, args.Error(1)
	}
	return products, args.Error(1)
}

func (m *MockProductRepository) FindByCategoryUUID(categoryUUID uuid.UUID, includeDeleted bool, isAvailable *bool) ([]models.Product, error) {
	args := m.Called(categoryUUID, includeDeleted, isAvailable)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	products, ok := args.Get(0).([]models.Product)
	if !ok {
		return nil, args.Error(1)
	}
	return products, args.Error(1)
}

func (m *MockProductRepository) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) ExistsBySlug(slug string) (bool, error) {
	args := m.Called(slug)
	return args.Bool(0), args.Error(1)
}

// Customization operations

func (m *MockProductRepository) CreateCustomization(customization *models.ProductCustomization) error {
	args := m.Called(customization)
	return args.Error(0)
}

func (m *MockProductRepository) FindCustomizationByUUID(uuid uuid.UUID) (*models.ProductCustomization, error) {
	args := m.Called(uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	customization, ok := args.Get(0).(*models.ProductCustomization)
	if !ok {
		return nil, args.Error(1)
	}
	return customization, args.Error(1)
}

func (m *MockProductRepository) UpdateCustomization(customization *models.ProductCustomization) error {
	args := m.Called(customization)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteCustomization(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) FindCustomizationsByProductID(productID uint) ([]models.ProductCustomization, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	customizations, ok := args.Get(0).([]models.ProductCustomization)
	if !ok {
		return nil, args.Error(1)
	}
	return customizations, args.Error(1)
}
