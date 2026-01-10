package routes

import (
	"github.com/carllix/matchaciee-backend/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupPaymentRoutes(
	app *fiber.App,
	paymentHandler *handlers.PaymentHandler,
) {
	api := app.Group("/api/v1")
	api.Post("/orders/:id/payment", paymentHandler.CreatePaymentToken)
	api.Post("/webhooks/midtrans", paymentHandler.HandleMidtransWebhook)
}
