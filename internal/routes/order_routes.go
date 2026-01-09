package routes

import (
	"github.com/carllix/matchaciee-backend/internal/handlers"
	"github.com/carllix/matchaciee-backend/internal/middleware"
	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func SetupOrderRoutes(
	app *fiber.App,
	orderHandler *handlers.OrderHandler,
	jwtUtil *utils.JWTUtil,
) {
	api := app.Group("/api/v1")
	orders := api.Group("/orders")

	// Public routes
	orders.Post("/guest", orderHandler.CreateGuestOrder)
	orders.Get("/track/:uuid", orderHandler.TrackGuestOrder)

	// Member routes
	orders.Post("/",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleMember, models.RoleAdmin),
		orderHandler.CreateOrder,
	)

	orders.Get("/me",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleMember),
		orderHandler.GetMyOrders,
	)

	orders.Get("/:id",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleMember, models.RoleAdmin, models.RoleBarista),
		orderHandler.GetOrder,
	)

	// Admin/Barista routes
	orders.Get("/",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin, models.RoleBarista),
		orderHandler.GetAllOrders,
	)

	orders.Get("/number/:number",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin),
		orderHandler.GetOrderByNumber,
	)

	orders.Put("/:id/status",
		middleware.AuthMiddleware(jwtUtil),
		middleware.RoleMiddleware(models.RoleAdmin, models.RoleBarista),
		orderHandler.UpdateOrderStatus,
	)
}
