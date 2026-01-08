package services

import (
	"errors"
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

func TestCategoryService_Create(t *testing.T) {
	t.Run("success - with slug provided", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		req := services.CreateCategoryRequest{
			Name: "Matcha Drinks",
			Slug: "matcha-drinks",
		}

		mockRepo.On("ExistsBySlug", "matcha-drinks").Return(false, nil)
		mockRepo.On("Create", mock.AnythingOfType("*models.Category")).Return(nil)

		result, err := service.Create(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Matcha Drinks", result.Name)
		assert.Equal(t, "matcha-drinks", result.Slug)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - auto-generate slug", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		req := services.CreateCategoryRequest{
			Name: "Iced Beverages",
		}

		mockRepo.On("ExistsBySlug", "iced-beverages").Return(false, nil)
		mockRepo.On("Create", mock.AnythingOfType("*models.Category")).Return(nil)

		result, err := service.Create(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Iced Beverages", result.Name)
		assert.Equal(t, "iced-beverages", result.Slug)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - slug already exists", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		req := services.CreateCategoryRequest{
			Name: "Matcha Drinks",
			Slug: "matcha-drinks",
		}

		mockRepo.On("ExistsBySlug", "matcha-drinks").Return(true, nil)

		result, err := service.Create(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, services.ErrCategorySlugExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository error on exists check", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		req := services.CreateCategoryRequest{
			Name: "Matcha Drinks",
		}

		mockRepo.On("ExistsBySlug", "matcha-drinks").Return(false, errors.New("database error"))

		result, err := service.Create(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository error on create", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		req := services.CreateCategoryRequest{
			Name: "Matcha Drinks",
		}

		mockRepo.On("ExistsBySlug", "matcha-drinks").Return(false, nil)
		mockRepo.On("Create", mock.AnythingOfType("*models.Category")).Return(errors.New("database error"))

		result, err := service.Create(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_GetByUUID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categoryUUID := uuid.New()
		category := &models.Category{
			ID:        1,
			UUID:      categoryUUID,
			Name:      "Matcha Drinks",
			Slug:      "matcha-drinks",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("FindByUUID", categoryUUID).Return(category, nil)

		result, err := service.GetByUUID(categoryUUID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, categoryUUID, result.ID)
		assert.Equal(t, "Matcha Drinks", result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categoryUUID := uuid.New()
		mockRepo.On("FindByUUID", categoryUUID).Return(nil, repositories.ErrCategoryNotFound)

		result, err := service.GetByUUID(categoryUUID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, services.ErrCategoryNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_GetBySlug(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		category := &models.Category{
			ID:        1,
			UUID:      uuid.New(),
			Name:      "Matcha Drinks",
			Slug:      "matcha-drinks",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("FindBySlug", "matcha-drinks").Return(category, nil)

		result, err := service.GetBySlug("matcha-drinks")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "matcha-drinks", result.Slug)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		mockRepo.On("FindBySlug", "non-existent").Return(nil, repositories.ErrCategoryNotFound)

		result, err := service.GetBySlug("non-existent")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, services.ErrCategoryNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_GetAll(t *testing.T) {
	t.Run("success - all categories", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categories := []models.Category{
			{ID: 1, UUID: uuid.New(), Name: "Category 1", Slug: "category-1", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, UUID: uuid.New(), Name: "Category 2", Slug: "category-2", IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		mockRepo.On("FindAll", (*bool)(nil)).Return(categories, nil)

		result, err := service.GetAll(false)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - active only", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categories := []models.Category{
			{ID: 1, UUID: uuid.New(), Name: "Category 1", Slug: "category-1", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		mockRepo.On("FindAll", mock.MatchedBy(func(isActive *bool) bool {
			return isActive != nil && *isActive == true
		})).Return(categories, nil)

		result, err := service.GetAll(true)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_Update(t *testing.T) {
	t.Run("success - update name", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categoryUUID := uuid.New()
		existingCategory := &models.Category{
			ID:        1,
			UUID:      categoryUUID,
			Name:      "Old Name",
			Slug:      "old-name",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		newName := "New Name"
		req := services.UpdateCategoryRequest{
			Name: &newName,
		}

		mockRepo.On("FindByUUID", categoryUUID).Return(existingCategory, nil)
		mockRepo.On("Update", mock.AnythingOfType("*models.Category")).Return(nil)

		result, err := service.Update(categoryUUID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - update slug", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categoryUUID := uuid.New()
		existingCategory := &models.Category{
			ID:        1,
			UUID:      categoryUUID,
			Name:      "Category",
			Slug:      "old-slug",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		newSlug := "new-slug"
		req := services.UpdateCategoryRequest{
			Slug: &newSlug,
		}

		mockRepo.On("FindByUUID", categoryUUID).Return(existingCategory, nil)
		mockRepo.On("ExistsBySlug", "new-slug").Return(false, nil)
		mockRepo.On("Update", mock.AnythingOfType("*models.Category")).Return(nil)

		result, err := service.Update(categoryUUID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "new-slug", result.Slug)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - slug already exists", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categoryUUID := uuid.New()
		existingCategory := &models.Category{
			ID:        1,
			UUID:      categoryUUID,
			Name:      "Category",
			Slug:      "old-slug",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		newSlug := "existing-slug"
		req := services.UpdateCategoryRequest{
			Slug: &newSlug,
		}

		mockRepo.On("FindByUUID", categoryUUID).Return(existingCategory, nil)
		mockRepo.On("ExistsBySlug", "existing-slug").Return(true, nil)

		result, err := service.Update(categoryUUID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, services.ErrCategorySlugExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categoryUUID := uuid.New()
		req := services.UpdateCategoryRequest{}

		mockRepo.On("FindByUUID", categoryUUID).Return(nil, repositories.ErrCategoryNotFound)

		result, err := service.Update(categoryUUID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, services.ErrCategoryNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categoryUUID := uuid.New()
		category := &models.Category{
			ID:        1,
			UUID:      categoryUUID,
			Name:      "Category",
			Slug:      "category",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("FindByUUID", categoryUUID).Return(category, nil)
		mockRepo.On("Delete", uint(1)).Return(nil)

		err := service.Delete(categoryUUID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categoryUUID := uuid.New()
		mockRepo.On("FindByUUID", categoryUUID).Return(nil, repositories.ErrCategoryNotFound)

		err := service.Delete(categoryUUID)

		assert.Error(t, err)
		assert.Equal(t, services.ErrCategoryNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockRepo := new(mocks.MockCategoryRepository)
		service := services.NewCategoryService(mockRepo)

		categoryUUID := uuid.New()
		category := &models.Category{
			ID:        1,
			UUID:      categoryUUID,
			Name:      "Category",
			Slug:      "category",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("FindByUUID", categoryUUID).Return(category, nil)
		mockRepo.On("Delete", uint(1)).Return(errors.New("database error"))

		err := service.Delete(categoryUUID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
