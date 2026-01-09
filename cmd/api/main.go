package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/carllix/matchaciee-backend/internal/config"
	"github.com/carllix/matchaciee-backend/internal/database"
	"github.com/carllix/matchaciee-backend/internal/handlers"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/carllix/matchaciee-backend/internal/routes"
	"github.com/carllix/matchaciee-backend/internal/services"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting %s in %s mode...", cfg.AppName, cfg.Env)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize app
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: errorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: joinOrigins(cfg.AllowedOrigins),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		dbConnected := database.IsConnected()
		status := "ok"
		if !dbConnected {
			status = "degraded"
		}

		return c.JSON(fiber.Map{
			"status":   status,
			"service":  cfg.AppName,
			"database": dbConnected,
		})
	})

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Matchaciee API",
			"version": "v1.0.0",
		})
	})

	// Initialize dependencies
	db := database.GetDB()
	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshTokenExpiry)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, refreshTokenRepo, jwtUtil)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	orderService := services.NewOrderService(orderRepo, productRepo, userRepo)
	paymentService := services.NewPaymentService(
		paymentRepo,
		orderRepo,
		cfg.MidtransServerKey,
		cfg.MidtransClientKey,
		cfg.MidtransEnvironment,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Setup routes
	routes.SetupAuthRoutes(app, authHandler, jwtUtil)
	routes.SetupProductRoutes(app, categoryHandler, productHandler, jwtUtil)
	routes.SetupOrderRoutes(app, orderHandler, jwtUtil)
	routes.SetupPaymentRoutes(app, paymentHandler)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server listening on port %s", cfg.AppPort)

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Block until we receive a signal
	<-quit
	log.Println("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	log.Println("Server stopped")
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}

func joinOrigins(origins []string) string {
	return strings.Join(origins, ", ")
}
