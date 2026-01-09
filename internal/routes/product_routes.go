package routes

import (
	"github.com/carllix/matchaciee-backend/internal/handlers"
	"github.com/carllix/matchaciee-backend/internal/middleware"
	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(
	app *fiber.App,
	categoryHandler *handlers.CategoryHandler,
	productHandler *handlers.ProductHandler,
	jwtUtil *utils.JWTUtil,
) {
	api := app.Group("/api/v1")

	// Category routes
	categories := api.Group("/categories")

	// Public routes
	categories.Get("/", categoryHandler.GetAllCategories)
	categories.Get("/:id", categoryHandler.GetCategory)
	categories.Get("/slug/:slug", categoryHandler.GetCategoryBySlug)

	// Admin routes
	categories.Post("/",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		categoryHandler.CreateCategory,
	)
	categories.Put("/:id",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		categoryHandler.UpdateCategory,
	)
	categories.Delete("/:id",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		categoryHandler.DeleteCategory,
	)

	// Product routes
	products := api.Group("/products")

	// Public routes
	products.Get("/", productHandler.GetAllProducts)
	products.Get("/:id", productHandler.GetProduct)
	products.Get("/slug/:slug", productHandler.GetProductBySlug)

	// Admin routes
	products.Post("/",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		productHandler.CreateProduct,
	)
	products.Put("/:id",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		productHandler.UpdateProduct,
	)
	products.Delete("/:id",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		productHandler.DeleteProduct,
	)
	products.Post("/:id/restore",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		productHandler.RestoreProduct,
	)

	// Product customization routes (Admin)
	products.Post("/:id/customizations",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		productHandler.AddProductCustomization,
	)
	products.Put("/customizations/:customizationId",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		productHandler.UpdateProductCustomization,
	)
	products.Delete("/customizations/:customizationId",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		productHandler.DeleteProductCustomization,
	)

	categories.Get("/:categoryId/products", productHandler.GetProductsByCategory)
}
