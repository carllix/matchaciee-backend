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

func TestProductService_Create(t *testing.T) {
	t.Run("success - with category", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		categoryUUID := uuid.New()
		category := &models.Category{
			ID:   1,
			UUID: categoryUUID,
			Name: "Drinks",
		}

		req := services.CreateProductRequest{
			Name:         "Matcha Latte",
			BasePrice:    45000,
			CategoryUUID: &categoryUUID,
		}

		mockProductRepo.On("ExistsBySlug", "matcha-latte").Return(false, nil)
		mockCategoryRepo.On("FindByUUID", categoryUUID).Return(category, nil)
		mockProductRepo.On("Create", mock.AnythingOfType("*models.Product")).Return(nil)

		result, err := service.Create(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Matcha Latte", result.Name)
		assert.Equal(t, "matcha-latte", result.Slug)
		assert.Equal(t, 45000.0, result.BasePrice)
		mockProductRepo.AssertExpectations(t)
		mockCategoryRepo.AssertExpectations(t)
	})

	t.Run("success - without category", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		req := services.CreateProductRequest{
			Name:      "Matcha Latte",
			BasePrice: 45000,
		}

		mockProductRepo.On("ExistsBySlug", "matcha-latte").Return(false, nil)
		mockProductRepo.On("Create", mock.AnythingOfType("*models.Product")).Return(nil)

		result, err := service.Create(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("success - with customizations", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		customizations := []services.CreateCustomizationRequest{
			{CustomizationType: "Size", OptionName: "Large", PriceModifier: 5000},
		}

		req := services.CreateProductRequest{
			Name:           "Matcha Latte",
			BasePrice:      45000,
			Customizations: customizations,
		}

		product := &models.Product{
			ID:   1,
			UUID: uuid.New(),
			Name: "Matcha Latte",
			Slug: "matcha-latte",
		}

		mockProductRepo.On("ExistsBySlug", "matcha-latte").Return(false, nil)
		mockProductRepo.On("Create", mock.AnythingOfType("*models.Product")).Return(nil)
		mockProductRepo.On("CreateCustomization", mock.AnythingOfType("*models.ProductCustomization")).Return(nil)
		mockProductRepo.On("FindByID", mock.AnythingOfType("uint")).Return(product, nil)

		result, err := service.Create(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - slug already exists", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		req := services.CreateProductRequest{
			Name:      "Matcha Latte",
			BasePrice: 45000,
		}

		mockProductRepo.On("ExistsBySlug", "matcha-latte").Return(true, nil)

		result, err := service.Create(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, services.ErrProductSlugExists, err)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		categoryUUID := uuid.New()
		req := services.CreateProductRequest{
			Name:         "Matcha Latte",
			BasePrice:    45000,
			CategoryUUID: &categoryUUID,
		}

		mockProductRepo.On("ExistsBySlug", "matcha-latte").Return(false, nil)
		mockCategoryRepo.On("FindByUUID", categoryUUID).Return(nil, repositories.ErrCategoryNotFound)

		result, err := service.Create(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, services.ErrCategoryNotFound, err)
		mockProductRepo.AssertExpectations(t)
		mockCategoryRepo.AssertExpectations(t)
	})
}

func TestProductService_GetByUUID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		product := &models.Product{
			ID:        1,
			UUID:      productUUID,
			Name:      "Matcha Latte",
			Slug:      "matcha-latte",
			BasePrice: 45000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockProductRepo.On("FindByUUID", productUUID).Return(product, nil)

		result, err := service.GetByUUID(productUUID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, productUUID, result.ID)
		assert.Equal(t, "Matcha Latte", result.Name)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - product not found", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		mockProductRepo.On("FindByUUID", productUUID).Return(nil, repositories.ErrProductNotFound)

		result, err := service.GetByUUID(productUUID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, services.ErrProductNotFound, err)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_GetAll(t *testing.T) {
	t.Run("success - all products", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		products := []models.Product{
			{ID: 1, UUID: uuid.New(), Name: "Product 1", Slug: "product-1", BasePrice: 10000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, UUID: uuid.New(), Name: "Product 2", Slug: "product-2", BasePrice: 20000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		mockProductRepo.On("FindAll", false, (*bool)(nil), (*uint)(nil)).Return(products, nil)

		result, err := service.GetAll(false, false, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_Update(t *testing.T) {
	t.Run("success - update name and price", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		existingProduct := &models.Product{
			ID:        1,
			UUID:      productUUID,
			Name:      "Old Name",
			Slug:      "old-name",
			BasePrice: 40000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		updatedProduct := &models.Product{
			ID:        1,
			UUID:      productUUID,
			Name:      "New Name",
			Slug:      "old-name",
			BasePrice: 50000,
			CreatedAt: existingProduct.CreatedAt,
			UpdatedAt: time.Now(),
		}

		newName := "New Name"
		newPrice := 50000.0
		req := services.UpdateProductRequest{
			Name:      &newName,
			BasePrice: &newPrice,
		}

		mockProductRepo.On("FindByUUID", productUUID).Return(existingProduct, nil)
		mockProductRepo.On("Update", mock.AnythingOfType("*models.Product")).Return(nil)
		mockProductRepo.On("FindByID", uint(1)).Return(updatedProduct, nil)

		result, err := service.Update(productUUID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, 50000.0, result.BasePrice)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - product not found", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		req := services.UpdateProductRequest{}

		mockProductRepo.On("FindByUUID", productUUID).Return(nil, repositories.ErrProductNotFound)

		result, err := service.Update(productUUID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, services.ErrProductNotFound, err)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_SoftDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		product := &models.Product{
			ID:   1,
			UUID: productUUID,
			Name: "Product",
		}

		mockProductRepo.On("FindByUUID", productUUID).Return(product, nil)
		mockProductRepo.On("SoftDelete", uint(1)).Return(nil)

		err := service.SoftDelete(productUUID)

		assert.NoError(t, err)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - product not found", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		mockProductRepo.On("FindByUUID", productUUID).Return(nil, repositories.ErrProductNotFound)

		err := service.SoftDelete(productUUID)

		assert.Error(t, err)
		assert.Equal(t, services.ErrProductNotFound, err)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_Restore(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		deletedAt := time.Now()
		product := &models.Product{
			ID:        1,
			UUID:      productUUID,
			Name:      "Product",
			DeletedAt: &deletedAt,
		}

		mockProductRepo.On("FindByUUIDIncludingDeleted", productUUID).Return(product, nil)
		mockProductRepo.On("Restore", uint(1)).Return(nil)

		err := service.Restore(productUUID)

		assert.NoError(t, err)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - product not deleted", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		product := &models.Product{
			ID:        1,
			UUID:      productUUID,
			Name:      "Product",
			DeletedAt: nil,
		}

		mockProductRepo.On("FindByUUIDIncludingDeleted", productUUID).Return(product, nil)

		err := service.Restore(productUUID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not deleted")
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - product not found", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		mockProductRepo.On("FindByUUIDIncludingDeleted", productUUID).Return(nil, repositories.ErrProductNotFound)

		err := service.Restore(productUUID)

		assert.Error(t, err)
		assert.Equal(t, services.ErrProductNotFound, err)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_AddCustomization(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		product := &models.Product{
			ID:   1,
			UUID: productUUID,
			Name: "Product",
		}

		req := services.CreateCustomizationRequest{
			CustomizationType: "Size",
			OptionName:        "Large",
			PriceModifier:     5000,
		}

		mockProductRepo.On("FindByUUID", productUUID).Return(product, nil)
		mockProductRepo.On("CreateCustomization", mock.AnythingOfType("*models.ProductCustomization")).Return(nil)

		result, err := service.AddCustomization(productUUID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Size", result.CustomizationType)
		assert.Equal(t, "Large", result.OptionName)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - product not found", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		productUUID := uuid.New()
		req := services.CreateCustomizationRequest{
			CustomizationType: "Size",
			OptionName:        "Large",
		}

		mockProductRepo.On("FindByUUID", productUUID).Return(nil, repositories.ErrProductNotFound)

		result, err := service.AddCustomization(productUUID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, services.ErrProductNotFound, err)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_UpdateCustomization(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		customizationUUID := uuid.New()
		customization := &models.ProductCustomization{
			ID:                1,
			UUID:              customizationUUID,
			ProductID:         1,
			CustomizationType: "Size",
			OptionName:        "Large",
			PriceModifier:     5000,
		}

		newOptionName := "Extra Large"
		req := services.UpdateCustomizationRequest{
			OptionName: &newOptionName,
		}

		mockProductRepo.On("FindCustomizationByUUID", customizationUUID).Return(customization, nil)
		mockProductRepo.On("UpdateCustomization", mock.AnythingOfType("*models.ProductCustomization")).Return(nil)

		result, err := service.UpdateCustomization(customizationUUID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Extra Large", result.OptionName)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - customization not found", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		customizationUUID := uuid.New()
		req := services.UpdateCustomizationRequest{}

		mockProductRepo.On("FindCustomizationByUUID", customizationUUID).Return(nil, errors.New("not found"))

		result, err := service.UpdateCustomization(customizationUUID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_DeleteCustomization(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		customizationUUID := uuid.New()
		customization := &models.ProductCustomization{
			ID:   1,
			UUID: customizationUUID,
		}

		mockProductRepo.On("FindCustomizationByUUID", customizationUUID).Return(customization, nil)
		mockProductRepo.On("DeleteCustomization", uint(1)).Return(nil)

		err := service.DeleteCustomization(customizationUUID)

		assert.NoError(t, err)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("error - customization not found", func(t *testing.T) {
		mockProductRepo := new(mocks.MockProductRepository)
		mockCategoryRepo := new(mocks.MockCategoryRepository)
		service := services.NewProductService(mockProductRepo, mockCategoryRepo)

		customizationUUID := uuid.New()
		mockProductRepo.On("FindCustomizationByUUID", customizationUUID).Return(nil, errors.New("not found"))

		err := service.DeleteCustomization(customizationUUID)

		assert.Error(t, err)
		mockProductRepo.AssertExpectations(t)
	})
}
