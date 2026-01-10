package handlers

import (
	"errors"
	"log"

	"github.com/carllix/matchaciee-backend/internal/services"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// POST /api/v1/orders/:id/payment
func (h *PaymentHandler) CreatePaymentToken(c *fiber.Ctx) error {
	// Get order UUID from URL params
	orderIDParam := c.Params("id")
	orderUUID, err := uuid.Parse(orderIDParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid order ID")
	}

	// Create payment token
	paymentToken, err := h.paymentService.CreatePaymentToken(orderUUID)
	if err != nil {
		if errors.Is(err, services.ErrOrderNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found")
		}
		if errors.Is(err, services.ErrPaymentAlreadyExists) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Payment already exists for this order")
		}
		log.Printf("Failed to create payment token: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create payment token")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"payment_id":   paymentToken.PaymentID,
		"token":        paymentToken.Token,
		"redirect_url": paymentToken.RedirectURL,
	})
}

// POST /api/v1/webhooks/midtrans
func (h *PaymentHandler) HandleMidtransWebhook(c *fiber.Ctx) error {
	// Parse webhook notification
	var notification services.MidtransNotification
	if err := c.BodyParser(&notification); err != nil {
		log.Printf("Failed to parse webhook: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Log webhook for debugging
	log.Printf("Received Midtrans webhook for order: %s, status: %s", notification.OrderID, notification.TransactionStatus)

	// Process webhook
	if err := h.paymentService.ProcessWebhookNotification(&notification); err != nil {
		if errors.Is(err, services.ErrInvalidSignature) {
			log.Printf("Invalid signature for order: %s", notification.OrderID)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid signature",
			})
		}
		if errors.Is(err, services.ErrPaymentNotFound) {
			log.Printf("Payment not found for order: %s", notification.OrderID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Payment not found",
			})
		}
		if errors.Is(err, services.ErrInvalidAmount) {
			log.Printf("Invalid amount for order: %s", notification.OrderID)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid amount",
			})
		}

		log.Printf("Failed to process webhook: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to process webhook",
		})
	}

	// Always return 200 OK to Midtrans
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}
