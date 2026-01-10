package services

import (
	"errors"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrProductSlugExists = errors.New("product slug already exists")
)

type CreateProductRequest struct {
	Name            string                       `json:"name" validate:"required,min=2,max=255"`
	Slug            string                       `json:"slug,omitempty" validate:"omitempty,min=2,max=255"`
	Description     *string                      `json:"description,omitempty"`
	CategoryUUID    *uuid.UUID                   `json:"category_id,omitempty"`
	BasePrice       float64                      `json:"base_price" validate:"required,gt=0"`
	PreparationTime *int                         `json:"preparation_time,omitempty" validate:"omitempty,gt=0"`
	DisplayOrder    int                          `json:"display_order,omitempty"`
	IsAvailable     *bool                        `json:"is_available,omitempty"`
	IsCustomizable  *bool                        `json:"is_customizable,omitempty"`
	ImageURL        *string                      `json:"image_url,omitempty" validate:"omitempty,url"`
	Customizations  []CreateCustomizationRequest `json:"customizations,omitempty"`
}

type UpdateProductRequest struct {
	Name            *string    `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Slug            *string    `json:"slug,omitempty" validate:"omitempty,min=2,max=255"`
	Description     *string    `json:"description,omitempty"`
	BasePrice       *float64   `json:"base_price,omitempty" validate:"omitempty,gt=0"`
	CategoryUUID    *uuid.UUID `json:"category_id,omitempty"`
	ImageURL        *string    `json:"image_url,omitempty" validate:"omitempty,url"`
	IsAvailable     *bool      `json:"is_available,omitempty"`
	IsCustomizable  *bool      `json:"is_customizable,omitempty"`
	PreparationTime *int       `json:"preparation_time,omitempty" validate:"omitempty,gt=0"`
	DisplayOrder    *int       `json:"display_order,omitempty"`
}

type CreateCustomizationRequest struct {
	CustomizationType string  `json:"customization_type" validate:"required,min=2,max=50"`
	OptionName        string  `json:"option_name" validate:"required,min=2,max=100"`
	PriceModifier     float64 `json:"price_modifier"`
	DisplayOrder      int     `json:"display_order,omitempty"`
}

type UpdateCustomizationRequest struct {
	CustomizationType *string  `json:"customization_type,omitempty" validate:"omitempty,min=2,max=50"`
	OptionName        *string  `json:"option_name,omitempty" validate:"omitempty,min=2,max=100"`
	PriceModifier     *float64 `json:"price_modifier,omitempty"`
	DisplayOrder      *int     `json:"display_order,omitempty"`
}

type ProductResponse struct {
	ID              uuid.UUID               `json:"id"`
	Name            string                  `json:"name"`
	Slug            string                  `json:"slug"`
	Description     *string                 `json:"description,omitempty"`
	Category        *CategoryResponse       `json:"category,omitempty"`
	BasePrice       float64                 `json:"base_price"`
	PreparationTime int                     `json:"preparation_time"`
	DisplayOrder    int                     `json:"display_order"`
	IsAvailable     bool                    `json:"is_available"`
	IsCustomizable  bool                    `json:"is_customizable"`
	ImageURL        *string                 `json:"image_url,omitempty"`
	DeletedAt       *string                 `json:"deleted_at,omitempty"`
	Customizations  []CustomizationResponse `json:"customizations,omitempty"`
	CreatedAt       string                  `json:"created_at"`
	UpdatedAt       string                  `json:"updated_at"`
}

type CustomizationResponse struct {
	ID                uuid.UUID `json:"id"`
	CustomizationType string    `json:"customization_type"`
	OptionName        string    `json:"option_name"`
	PriceModifier     float64   `json:"price_modifier"`
	DisplayOrder      int       `json:"display_order"`
	CreatedAt         string    `json:"created_at"`
}

type ProductService interface {
	Create(req CreateProductRequest) (*ProductResponse, error)
	GetByUUID(uuid uuid.UUID) (*ProductResponse, error)
	GetBySlug(slug string) (*ProductResponse, error)
	GetAll(includeDeleted bool, availableOnly bool, categoryUUID *uuid.UUID) ([]ProductResponse, error)
	Update(uuid uuid.UUID, req UpdateProductRequest) (*ProductResponse, error)
	SoftDelete(uuid uuid.UUID) error
	Restore(uuid uuid.UUID) error

	AddCustomization(productUUID uuid.UUID, req CreateCustomizationRequest) (*CustomizationResponse, error)
	UpdateCustomization(customizationUUID uuid.UUID, req UpdateCustomizationRequest) (*CustomizationResponse, error)
	DeleteCustomization(customizationUUID uuid.UUID) error
}

type productService struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductService(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) Create(req CreateProductRequest) (*ProductResponse, error) {
	productSlug := req.Slug
	if productSlug == "" {
		productSlug = slug.Make(req.Name)
	} else {
		productSlug = slug.Make(productSlug)
	}

	exists, err := s.productRepo.ExistsBySlug(productSlug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrProductSlugExists
	}

	var categoryID *uint
	if req.CategoryUUID != nil {
		var category *models.Category
		category, err = s.categoryRepo.FindByUUID(*req.CategoryUUID)
		if err != nil {
			if errors.Is(err, repositories.ErrCategoryNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
		categoryID = &category.ID
	}

	isAvailable := true
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	isCustomizable := false
	if req.IsCustomizable != nil {
		isCustomizable = *req.IsCustomizable
	}

	preparationTime := 5
	if req.PreparationTime != nil {
		preparationTime = *req.PreparationTime
	}

	product := &models.Product{
		Name:            req.Name,
		Slug:            productSlug,
		Description:     req.Description,
		BasePrice:       req.BasePrice,
		CategoryID:      categoryID,
		ImageURL:        req.ImageURL,
		IsAvailable:     isAvailable,
		IsCustomizable:  isCustomizable,
		PreparationTime: preparationTime,
		DisplayOrder:    req.DisplayOrder,
	}

	err = s.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	// Add customizations if provided
	if len(req.Customizations) > 0 {
		for _, customizationReq := range req.Customizations {
			customization := &models.ProductCustomization{
				ProductID:         product.ID,
				CustomizationType: customizationReq.CustomizationType,
				OptionName:        customizationReq.OptionName,
				PriceModifier:     customizationReq.PriceModifier,
				DisplayOrder:      customizationReq.DisplayOrder,
			}
			err = s.productRepo.CreateCustomization(customization)
			if err != nil {
				return nil, err
			}
		}
		// Reload product with customizations
		product, err = s.productRepo.FindByID(product.ID)
		if err != nil {
			return nil, err
		}
	}

	return s.toProductResponse(product), nil
}

func (s *productService) GetByUUID(productUUID uuid.UUID) (*ProductResponse, error) {
	product, err := s.productRepo.FindByUUID(productUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *productService) GetBySlug(productSlug string) (*ProductResponse, error) {
	product, err := s.productRepo.FindBySlug(productSlug)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *productService) GetAll(includeDeleted bool, availableOnly bool, categoryUUID *uuid.UUID) ([]ProductResponse, error) {
	var isAvailable *bool
	if availableOnly {
		available := true
		isAvailable = &available
	}

	var categoryID *uint
	if categoryUUID != nil {
		category, err := s.categoryRepo.FindByUUID(*categoryUUID)
		if err != nil {
			if errors.Is(err, repositories.ErrCategoryNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
		categoryID = &category.ID
	}

	products, err := s.productRepo.FindAll(includeDeleted, isAvailable, categoryID)
	if err != nil {
		return nil, err
	}

	responses := make([]ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *s.toProductResponse(&product)
	}

	return responses, nil
}

func (s *productService) Update(productUUID uuid.UUID, req UpdateProductRequest) (*ProductResponse, error) {
	// Find existing product
	product, err := s.productRepo.FindByUUID(productUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		product.Name = *req.Name
	}

	if req.Slug != nil {
		newSlug := slug.Make(*req.Slug)
		if newSlug != product.Slug {
			var exists bool
			exists, err = s.productRepo.ExistsBySlug(newSlug)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrProductSlugExists
			}
			product.Slug = newSlug
		}
	}

	if req.Description != nil {
		product.Description = req.Description
	}

	if req.BasePrice != nil {
		product.BasePrice = *req.BasePrice
	}

	if req.CategoryUUID != nil {
		var category *models.Category
		category, err = s.categoryRepo.FindByUUID(*req.CategoryUUID)
		if err != nil {
			if errors.Is(err, repositories.ErrCategoryNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
		product.CategoryID = &category.ID
	}

	if req.ImageURL != nil {
		product.ImageURL = req.ImageURL
	}

	if req.IsAvailable != nil {
		product.IsAvailable = *req.IsAvailable
	}

	if req.IsCustomizable != nil {
		product.IsCustomizable = *req.IsCustomizable
	}

	if req.PreparationTime != nil {
		product.PreparationTime = *req.PreparationTime
	}

	if req.DisplayOrder != nil {
		product.DisplayOrder = *req.DisplayOrder
	}

	err = s.productRepo.Update(product)
	if err != nil {
		return nil, err
	}

	// Reload to get updated data with relations
	product, err = s.productRepo.FindByID(product.ID)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *productService) SoftDelete(productUUID uuid.UUID) error {
	product, err := s.productRepo.FindByUUID(productUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	return s.productRepo.SoftDelete(product.ID)
}

func (s *productService) Restore(productUUID uuid.UUID) error {
	// Find the product including deleted ones
	product, err := s.productRepo.FindByUUIDIncludingDeleted(productUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	// Check if product is actually deleted
	if product.DeletedAt == nil {
		return errors.New("product is not deleted")
	}

	return s.productRepo.Restore(product.ID)
}

func (s *productService) AddCustomization(productUUID uuid.UUID, req CreateCustomizationRequest) (*CustomizationResponse, error) {
	product, err := s.productRepo.FindByUUID(productUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if !product.IsCustomizable {
		return nil, ErrProductNotCustomizable
	}

	customization := &models.ProductCustomization{
		ProductID:         product.ID,
		CustomizationType: req.CustomizationType,
		OptionName:        req.OptionName,
		PriceModifier:     req.PriceModifier,
		DisplayOrder:      req.DisplayOrder,
	}

	err = s.productRepo.CreateCustomization(customization)
	if err != nil {
		return nil, err
	}

	return s.toCustomizationResponse(customization), nil
}

func (s *productService) UpdateCustomization(customizationUUID uuid.UUID, req UpdateCustomizationRequest) (*CustomizationResponse, error) {
	// Find existing customization
	customization, err := s.productRepo.FindCustomizationByUUID(customizationUUID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.CustomizationType != nil {
		customization.CustomizationType = *req.CustomizationType
	}

	if req.OptionName != nil {
		customization.OptionName = *req.OptionName
	}

	if req.PriceModifier != nil {
		customization.PriceModifier = *req.PriceModifier
	}

	if req.DisplayOrder != nil {
		customization.DisplayOrder = *req.DisplayOrder
	}

	err = s.productRepo.UpdateCustomization(customization)
	if err != nil {
		return nil, err
	}

	return s.toCustomizationResponse(customization), nil
}

func (s *productService) DeleteCustomization(customizationUUID uuid.UUID) error {
	// Find existing customization
	customization, err := s.productRepo.FindCustomizationByUUID(customizationUUID)
	if err != nil {
		return err
	}

	return s.productRepo.DeleteCustomization(customization.ID)
}

func (s *productService) toProductResponse(product *models.Product) *ProductResponse {
	response := &ProductResponse{
		ID:              product.UUID,
		Name:            product.Name,
		Slug:            product.Slug,
		Description:     product.Description,
		BasePrice:       product.BasePrice,
		ImageURL:        product.ImageURL,
		IsAvailable:     product.IsAvailable,
		IsCustomizable:  product.IsCustomizable,
		PreparationTime: product.PreparationTime,
		DisplayOrder:    product.DisplayOrder,
		CreatedAt:       product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if product.DeletedAt != nil {
		deletedAtStr := product.DeletedAt.Format("2006-01-02T15:04:05Z07:00")
		response.DeletedAt = &deletedAtStr
	}

	if product.Category != nil {
		response.Category = &CategoryResponse{
			ID:           product.Category.UUID,
			Name:         product.Category.Name,
			Slug:         product.Category.Slug,
			Description:  product.Category.Description,
			ImageURL:     product.Category.ImageURL,
			DisplayOrder: product.Category.DisplayOrder,
			IsActive:     product.Category.IsActive,
			CreatedAt:    product.Category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    product.Category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	if len(product.Customizations) > 0 {
		response.Customizations = make([]CustomizationResponse, len(product.Customizations))
		for i, customization := range product.Customizations {
			response.Customizations[i] = *s.toCustomizationResponse(&customization)
		}
	}

	return response
}

func (s *productService) toCustomizationResponse(customization *models.ProductCustomization) *CustomizationResponse {
	return &CustomizationResponse{
		ID:                customization.UUID,
		CustomizationType: customization.CustomizationType,
		OptionName:        customization.OptionName,
		PriceModifier:     customization.PriceModifier,
		DisplayOrder:      customization.DisplayOrder,
		CreatedAt:         customization.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
