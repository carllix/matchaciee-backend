package routes

import (
	"github.com/carllix/matchaciee-backend/internal/handlers"
	"github.com/carllix/matchaciee-backend/internal/middleware"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, authHandler *handlers.AuthHandler, jwtUtil *utils.JWTUtil) {
	auth := app.Group("/api/v1/auth")

	// Public routes
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", authHandler.Logout)

	// Protected routes
	auth.Get("/me", middleware.AuthMiddleware(jwtUtil), authHandler.GetMe)
}
