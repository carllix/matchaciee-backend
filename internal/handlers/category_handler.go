package handlers

import (
	"errors"

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

// POST /api/v1/categories
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

// GET /api/v1/categories/:id
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

// GET /api/v1/categories/slug/:slug
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

// GET /api/v1/categories
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

// PUT /api/v1/categories/:id
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

// DELETE /api/v1/categories/:id
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
