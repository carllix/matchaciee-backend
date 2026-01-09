package mocks

import (
	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) FindByID(id uint) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	category, ok := args.Get(0).(*models.Category)
	if !ok {
		return nil, args.Error(1)
	}
	return category, args.Error(1)
}

func (m *MockCategoryRepository) FindByUUID(uuid uuid.UUID) (*models.Category, error) {
	args := m.Called(uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	category, ok := args.Get(0).(*models.Category)
	if !ok {
		return nil, args.Error(1)
	}
	return category, args.Error(1)
}

func (m *MockCategoryRepository) FindBySlug(slug string) (*models.Category, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	category, ok := args.Get(0).(*models.Category)
	if !ok {
		return nil, args.Error(1)
	}
	return category, args.Error(1)
}

func (m *MockCategoryRepository) FindAll(isActive *bool) ([]models.Category, error) {
	args := m.Called(isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	categories, ok := args.Get(0).([]models.Category)
	if !ok {
		return nil, args.Error(1)
	}
	return categories, args.Error(1)
}

func (m *MockCategoryRepository) Update(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCategoryRepository) ExistsBySlug(slug string) (bool, error) {
	args := m.Called(slug)
	return args.Bool(0), args.Error(1)
}

func (m *MockCategoryRepository) ExistsByName(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}
