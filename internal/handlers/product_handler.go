package handlers

import (
	"errors"

	"github.com/carllix/matchaciee-backend/internal/services"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req services.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Create product
	product, err := h.productService.Create(req)
	if err != nil {
		if errors.Is(err, services.ErrProductSlugExists) {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Product slug already exists")
		}
		if errors.Is(err, services.ErrCategoryNotFound) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create product")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, product)
}

// GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Parse as UUID
	productUUID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID format")
	}

	product, err := h.productService.GetByUUID(productUUID)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Product not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get product")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, product)
}

// GET /api/v1/products/slug/:slug
func (h *ProductHandler) GetProductBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	product, err := h.productService.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Product not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get product")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, product)
}

// GET /api/v1/products
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	includeDeleted := c.QueryBool("include_deleted", false)
	availableOnly := c.QueryBool("available_only", false)

	var categoryUUID *uuid.UUID
	categoryParam := c.Query("category_id")
	if categoryParam != "" {
		parsed, err := uuid.Parse(categoryParam)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category_id format")
		}
		categoryUUID = &parsed
	}

	products, err := h.productService.GetAll(includeDeleted, availableOnly, categoryUUID)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get products")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"products": products,
		"count":    len(products),
	})
}

// PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Parse UUID
	productUUID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID format")
	}

	var req services.UpdateProductRequest
	if err = c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Update product
	product, err := h.productService.Update(productUUID, req)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Product not found")
		}
		if errors.Is(err, services.ErrProductSlugExists) {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Product slug already exists")
		}
		if errors.Is(err, services.ErrCategoryNotFound) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update product")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, product)
}

// DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Parse UUID
	productUUID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID format")
	}

	// Soft delete product
	err = h.productService.SoftDelete(productUUID)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Product not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete product")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Product deleted successfully",
	})
}

// POST /api/v1/products/:id/restore
func (h *ProductHandler) RestoreProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Parse UUID
	productUUID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID format")
	}

	// Restore product
	err = h.productService.Restore(productUUID)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Product not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to restore product")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Product restored successfully",
	})
}

// POST /api/v1/products/:id/customizations
func (h *ProductHandler) AddProductCustomization(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Parse UUID
	productUUID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID format")
	}

	var req services.CreateCustomizationRequest
	if err = c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Add customization
	customization, err := h.productService.AddCustomization(productUUID, req)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Product not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add customization")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, customization)
}

// PUT /api/v1/products/customizations/:customizationId
func (h *ProductHandler) UpdateProductCustomization(c *fiber.Ctx) error {
	customizationParam := c.Params("customizationId")

	// Parse UUID
	customizationUUID, err := uuid.Parse(customizationParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid customization ID format")
	}

	var req services.UpdateCustomizationRequest
	if err = c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Update customization
	customization, err := h.productService.UpdateCustomization(customizationUUID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update customization")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, customization)
}

// DELETE /api/v1/products/customizations/:customizationId
func (h *ProductHandler) DeleteProductCustomization(c *fiber.Ctx) error {
	customizationParam := c.Params("customizationId")

	// Parse UUID
	customizationUUID, err := uuid.Parse(customizationParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid customization ID format")
	}

	// Delete customization
	err = h.productService.DeleteCustomization(customizationUUID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete customization")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Customization deleted successfully",
	})
}
