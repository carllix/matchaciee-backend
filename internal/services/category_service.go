package services

import (
	"errors"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

var (
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategorySlugExists = errors.New("category slug already exists")
)

type CreateCategoryRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=100"`
	Slug         string  `json:"slug,omitempty" validate:"omitempty,min=2,max=100"`
	Description  *string `json:"description,omitempty"`
	DisplayOrder int     `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
	ImageURL     *string `json:"image_url,omitempty" validate:"omitempty,url"`
}

type UpdateCategoryRequest struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Slug         *string `json:"slug,omitempty" validate:"omitempty,min=2,max=100"`
	Description  *string `json:"description,omitempty"`
	ImageURL     *string `json:"image_url,omitempty" validate:"omitempty,url"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

type CategoryResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description,omitempty"`
	ImageURL     *string   `json:"image_url,omitempty"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type CategoryService interface {
	Create(req CreateCategoryRequest) (*CategoryResponse, error)
	GetByUUID(uuid uuid.UUID) (*CategoryResponse, error)
	GetBySlug(slug string) (*CategoryResponse, error)
	GetAll(activeOnly bool) ([]CategoryResponse, error)
	Update(uuid uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error)
	Delete(uuid uuid.UUID) error
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) Create(req CreateCategoryRequest) (*CategoryResponse, error) {
	categorySlug := req.Slug
	if categorySlug == "" {
		categorySlug = slug.Make(req.Name)
	} else {
		categorySlug = slug.Make(categorySlug)
	}

	exists, err := s.categoryRepo.ExistsBySlug(categorySlug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCategorySlugExists
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	category := &models.Category{
		Name:         req.Name,
		Slug:         categorySlug,
		Description:  req.Description,
		ImageURL:     req.ImageURL,
		DisplayOrder: req.DisplayOrder,
		IsActive:     isActive,
	}

	err = s.categoryRepo.Create(category)
	if err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *categoryService) GetByUUID(uuid uuid.UUID) (*CategoryResponse, error) {
	category, err := s.categoryRepo.FindByUUID(uuid)
	if err != nil {
		if errors.Is(err, repositories.ErrCategoryNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *categoryService) GetBySlug(categorySlug string) (*CategoryResponse, error) {
	category, err := s.categoryRepo.FindBySlug(categorySlug)
	if err != nil {
		if errors.Is(err, repositories.ErrCategoryNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *categoryService) GetAll(activeOnly bool) ([]CategoryResponse, error) {
	var isActive *bool
	if activeOnly {
		active := true
		isActive = &active
	}

	categories, err := s.categoryRepo.FindAll(isActive)
	if err != nil {
		return nil, err
	}

	responses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *s.toCategoryResponse(&category)
	}

	return responses, nil
}

func (s *categoryService) Update(categoryUUID uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error) {
	category, err := s.categoryRepo.FindByUUID(categoryUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrCategoryNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}

	if req.Slug != nil {
		newSlug := slug.Make(*req.Slug)
		if newSlug != category.Slug {
			var exists bool
			exists, err = s.categoryRepo.ExistsBySlug(newSlug)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrCategorySlugExists
			}
			category.Slug = newSlug
		}
	}

	if req.Description != nil {
		category.Description = req.Description
	}

	if req.ImageURL != nil {
		category.ImageURL = req.ImageURL
	}

	if req.DisplayOrder != nil {
		category.DisplayOrder = *req.DisplayOrder
	}

	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	err = s.categoryRepo.Update(category)
	if err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *categoryService) Delete(categoryUUID uuid.UUID) error {
	category, err := s.categoryRepo.FindByUUID(categoryUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrCategoryNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	return s.categoryRepo.Delete(category.ID)
}

func (s *categoryService) toCategoryResponse(category *models.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:           category.UUID,
		Name:         category.Name,
		Slug:         category.Slug,
		Description:  category.Description,
		ImageURL:     category.ImageURL,
		DisplayOrder: category.DisplayOrder,
		IsActive:     category.IsActive,
		CreatedAt:    category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
