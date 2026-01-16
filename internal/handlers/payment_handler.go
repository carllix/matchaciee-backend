package handlers

import (
	"errors"
	"log"

	_ "github.com/carllix/matchaciee-backend/docs"
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

// CreatePaymentToken godoc
// @Summary Create payment token
// @Description Create a Midtrans payment token for an order. Returns a redirect URL and token for Snap payment.
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Order UUID"
// @Success 200 {object} docs.PaymentSuccessResponse "Payment token created successfully"
// @Failure 400 {object} docs.SwaggerErrorResponse "Invalid order ID or payment already exists"
// @Failure 404 {object} docs.SwaggerErrorResponse "Order not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /orders/{id}/payment [post]
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

// HandleMidtransWebhook godoc
// @Summary Handle Midtrans webhook
// @Description Process payment notification from Midtrans. This endpoint is called by Midtrans when payment status changes.
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param notification body docs.MidtransWebhookRequest true "Midtrans notification payload"
// @Success 200 {object} docs.WebhookSuccessResponse "Webhook processed successfully"
// @Failure 400 {object} docs.WebhookErrorResponse "Invalid request body or invalid amount"
// @Failure 401 {object} docs.WebhookErrorResponse "Invalid signature"
// @Failure 404 {object} docs.WebhookErrorResponse "Payment not found"
// @Failure 500 {object} docs.WebhookErrorResponse "Internal server error"
// @Router /webhooks/midtrans [post]
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
