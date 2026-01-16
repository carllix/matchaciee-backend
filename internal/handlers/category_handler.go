package handlers

import (
	"errors"

	_ "github.com/carllix/matchaciee-backend/docs"
	"github.com/carllix/matchaciee-backend/internal/services"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new product category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body docs.CreateCategoryRequest true "Category details"
// @Success 201 {object} docs.CategorySuccessResponse "Category created successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 409 {object} docs.SwaggerErrorResponse "Category slug already exists"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req services.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Create category
	category, err := h.categoryService.Create(req)
	if err != nil {
		if errors.Is(err, services.ErrCategorySlugExists) {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Category slug already exists")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create category")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, category)
}

// GetCategory godoc
// @Summary Get category by ID
// @Description Get a single category by its UUID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category UUID"
// @Success 200 {object} docs.CategorySuccessResponse "Category retrieved successfully"
// @Failure 400 {object} docs.SwaggerErrorResponse "Invalid category ID format"
// @Failure 404 {object} docs.SwaggerErrorResponse "Category not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Parse as UUID
	categoryUUID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID format")
	}

	category, err := h.categoryService.GetByUUID(categoryUUID)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get category")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, category)
}

// GetCategoryBySlug godoc
// @Summary Get category by slug
// @Description Get a single category by its URL-friendly slug
// @Tags Categories
// @Accept json
// @Produce json
// @Param slug path string true "Category slug"
// @Success 200 {object} docs.CategorySuccessResponse "Category retrieved successfully"
// @Failure 404 {object} docs.SwaggerErrorResponse "Category not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	category, err := h.categoryService.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get category")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, category)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Get a list of all categories with optional filtering
// @Tags Categories
// @Accept json
// @Produce json
// @Param active_only query boolean false "Filter to show only active categories"
// @Success 200 {object} docs.CategoriesSuccessResponse "Categories retrieved successfully"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	activeOnly := c.QueryBool("active_only", false)

	categories, err := h.categoryService.GetAll(activeOnly)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get categories")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"categories": categories,
		"count":      len(categories),
	})
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category by its UUID (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category UUID"
// @Param request body docs.UpdateCategoryRequest true "Category update details"
// @Success 200 {object} docs.CategorySuccessResponse "Category updated successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error or invalid ID format"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} docs.SwaggerErrorResponse "Category not found"
// @Failure 409 {object} docs.SwaggerErrorResponse "Category slug already exists"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Parse UUID
	categoryUUID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID format")
	}

	var req services.UpdateCategoryRequest
	if err = c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Update category
	category, err := h.categoryService.Update(categoryUUID, req)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
		}
		if errors.Is(err, services.ErrCategorySlugExists) {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Category slug already exists")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update category")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, category)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by its UUID (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category UUID"
// @Success 200 {object} docs.MessageSuccessResponse "Category deleted successfully"
// @Failure 400 {object} docs.SwaggerErrorResponse "Invalid category ID format"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} docs.SwaggerErrorResponse "Category not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Parse UUID
	categoryUUID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID format")
	}

	// Delete category
	err = h.categoryService.Delete(categoryUUID)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete category")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Category deleted successfully",
	})
}
