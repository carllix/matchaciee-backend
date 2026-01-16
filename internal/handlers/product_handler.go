package handlers

import (
	"errors"

	_ "github.com/carllix/matchaciee-backend/docs"
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

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with optional customizations (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body docs.CreateProductRequest true "Product details"
// @Success 201 {object} docs.ProductSuccessResponse "Product created successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error or category not found"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 409 {object} docs.SwaggerErrorResponse "Product slug already exists"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /products [post]
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

// GetProduct godoc
// @Summary Get product by ID
// @Description Get a single product by its UUID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product UUID"
// @Success 200 {object} docs.ProductSuccessResponse "Product retrieved successfully"
// @Failure 400 {object} docs.SwaggerErrorResponse "Invalid product ID format"
// @Failure 404 {object} docs.SwaggerErrorResponse "Product not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /products/{id} [get]
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

// GetProductBySlug godoc
// @Summary Get product by slug
// @Description Get a single product by its URL-friendly slug
// @Tags Products
// @Accept json
// @Produce json
// @Param slug path string true "Product slug"
// @Success 200 {object} docs.ProductSuccessResponse "Product retrieved successfully"
// @Failure 404 {object} docs.SwaggerErrorResponse "Product not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /products/slug/{slug} [get]
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

// GetAllProducts godoc
// @Summary Get all products
// @Description Get a list of all products with optional filtering
// @Tags Products
// @Accept json
// @Produce json
// @Param include_deleted query boolean false "Include soft-deleted products"
// @Param available_only query boolean false "Filter to show only available products"
// @Param category_id query string false "Filter by category UUID"
// @Success 200 {object} docs.ProductsSuccessResponse "Products retrieved successfully"
// @Failure 400 {object} docs.SwaggerErrorResponse "Invalid category_id format or category not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /products [get]
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

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product by its UUID (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product UUID"
// @Param request body docs.UpdateProductRequest true "Product update details"
// @Success 200 {object} docs.ProductSuccessResponse "Product updated successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error, invalid ID format, or category not found"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} docs.SwaggerErrorResponse "Product not found"
// @Failure 409 {object} docs.SwaggerErrorResponse "Product slug already exists"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /products/{id} [put]
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

// DeleteProduct godoc
// @Summary Soft delete a product
// @Description Soft delete a product by its UUID (Admin only). Product will be hidden but preserved for historical orders.
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product UUID"
// @Success 200 {object} docs.MessageSuccessResponse "Product deleted successfully"
// @Failure 400 {object} docs.SwaggerErrorResponse "Invalid product ID format"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} docs.SwaggerErrorResponse "Product not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /products/{id} [delete]
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

// RestoreProduct godoc
// @Summary Restore a soft-deleted product
// @Description Restore a previously soft-deleted product by its UUID (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product UUID"
// @Success 200 {object} docs.MessageSuccessResponse "Product restored successfully"
// @Failure 400 {object} docs.SwaggerErrorResponse "Invalid product ID format"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} docs.SwaggerErrorResponse "Product not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /products/{id}/restore [post]
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

// AddProductCustomization godoc
// @Summary Add customization to a product
// @Description Add a new customization option to a product (Admin only). Product must be marked as customizable.
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product UUID"
// @Param request body docs.CreateCustomizationRequest true "Customization details"
// @Success 201 {object} docs.CustomizationSuccessResponse "Customization added successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error, invalid ID format, or product is not customizable"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} docs.SwaggerErrorResponse "Product not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /products/{id}/customizations [post]
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
		if errors.Is(err, services.ErrProductNotCustomizable) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Product is not customizable")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add customization")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, customization)
}

// UpdateProductCustomization godoc
// @Summary Update a product customization
// @Description Update an existing product customization by its UUID (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param customizationId path string true "Customization UUID"
// @Param request body docs.UpdateCustomizationRequest true "Customization update details"
// @Success 200 {object} docs.CustomizationSuccessResponse "Customization updated successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error or invalid ID format"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /products/customizations/{customizationId} [put]
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

// DeleteProductCustomization godoc
// @Summary Delete a product customization
// @Description Delete a product customization by its UUID (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param customizationId path string true "Customization UUID"
// @Success 200 {object} docs.MessageSuccessResponse "Customization deleted successfully"
// @Failure 400 {object} docs.SwaggerErrorResponse "Invalid customization ID format"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 403 {object} docs.SwaggerErrorResponse "Forbidden - Admin only"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /products/customizations/{customizationId} [delete]
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
