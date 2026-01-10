package services

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/datatypes"
)

var (
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrInvalidSignature     = errors.New("invalid signature")
	ErrPaymentAlreadyExists = errors.New("payment already processed")
	ErrInvalidAmount        = errors.New("invalid payment amount")
)

type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

type MidtransNotification struct {
	TransactionTime   string  `json:"transaction_time"`
	TransactionStatus string  `json:"transaction_status"`
	TransactionID     string  `json:"transaction_id"`
	StatusMessage     string  `json:"status_message"`
	StatusCode        string  `json:"status_code"`
	SignatureKey      string  `json:"signature_key"`
	PaymentType       string  `json:"payment_type"`
	OrderID           string  `json:"order_id"`
	MerchantID        string  `json:"merchant_id"`
	GrossAmount       string  `json:"gross_amount"`
	FraudStatus       string  `json:"fraud_status"`
	Currency          string  `json:"currency"`
	SettlementTime    *string `json:"settlement_time,omitempty"`
}

type PaymentTokenResponse struct {
	PaymentID   uuid.UUID `json:"payment_id"`
	Token       string    `json:"token"`
	RedirectURL string    `json:"redirect_url"`
}

type PaymentService interface {
	CreatePaymentToken(orderUUID uuid.UUID) (*PaymentTokenResponse, error)
	ProcessWebhookNotification(notification *MidtransNotification) error
	VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool
}

type paymentService struct {
	paymentRepo repositories.PaymentRepository
	orderRepo   repositories.OrderRepository
	snapClient  snap.Client
	serverKey   string
}

func NewPaymentService(
	paymentRepo repositories.PaymentRepository,
	orderRepo repositories.OrderRepository,
	serverKey string,
	clientKey string,
	environment string,
) PaymentService {
	// Initialize Snap client
	var snapClient snap.Client
	if environment == "production" {
		snapClient.New(serverKey, midtrans.Production)
	} else {
		snapClient.New(serverKey, midtrans.Sandbox)
	}

	return &paymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		snapClient:  snapClient,
		serverKey:   serverKey,
	}
}

func (s *paymentService) CreatePaymentToken(orderUUID uuid.UUID) (*PaymentTokenResponse, error) {
	// Get order details
	order, err := s.orderRepo.FindByUUID(orderUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Validate order status
	if order.Status != models.OrderStatusPending {
		return nil, fmt.Errorf("order must be in pending status to create payment")
	}

	// Generate unique Midtrans order ID
	midtransOrderID := fmt.Sprintf("%s-%d", order.OrderNumber, time.Now().Unix())

	// Check if payment already exists for this order
	existingPayments, err := s.paymentRepo.FindByOrderID(order.ID)
	if err != nil {
		return nil, err
	}

	for _, p := range existingPayments {
		if p.TransactionStatus != nil && *p.TransactionStatus == models.TransactionStatusSettlement {
			return nil, ErrPaymentAlreadyExists
		}
	}

	// Prepare Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  midtransOrderID,
			GrossAmt: int64(order.Total),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.CustomerName,
			Email: func() string {
				if order.User != nil {
					return order.User.Email
				}
				return ""
			}(),
		},
	}

	// Add item details
	var items []midtrans.ItemDetails
	for _, item := range order.Items {
		quantity := item.Quantity
		if quantity > 2147483647 || quantity < 0 {
			return nil, fmt.Errorf("invalid quantity for item %s", item.ProductName)
		}

		items = append(items, midtrans.ItemDetails{
			ID:    item.UUID.String(),
			Name:  item.ProductName,
			Price: int64(item.UnitPrice),
			Qty:   int32(quantity),
		})
	}
	req.Items = &items

	// Create Snap transaction
	snapResp, midtransErr := s.snapClient.CreateTransaction(req)
	if midtransErr != nil {
		log.Printf("Failed to create Snap transaction: %v", midtransErr)
		return nil, fmt.Errorf("failed to create payment: %w", midtransErr)
	}

	// Create payment record
	payment := &models.Payment{
		OrderID:         order.ID,
		MidtransOrderID: midtransOrderID,
		GrossAmount:     order.Total,
		PaymentMetadata: datatypes.JSON("{}"),
	}

	err = s.paymentRepo.Create(payment)
	if err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return &PaymentTokenResponse{
		PaymentID:   payment.UUID,
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

func (s *paymentService) ProcessWebhookNotification(notification *MidtransNotification) error {
	// Verify signature
	if !s.VerifySignature(notification.OrderID, notification.StatusCode, notification.GrossAmount, notification.SignatureKey) {
		log.Printf("Invalid signature for order: %s", notification.OrderID)
		return ErrInvalidSignature
	}

	// Find payment by Midtrans order ID
	payment, err := s.paymentRepo.FindByMidtransOrderID(notification.OrderID)
	if err != nil {
		if errors.Is(err, repositories.ErrPaymentNotFound) {
			return ErrPaymentNotFound
		}
		return err
	}

	// Verify gross amount matches
	grossAmount, err := strconv.ParseFloat(notification.GrossAmount, 64)
	if err != nil {
		return ErrInvalidAmount
	}

	if grossAmount != payment.GrossAmount {
		log.Printf("Amount mismatch for order %s: expected %.2f, got %.2f", notification.OrderID, payment.GrossAmount, grossAmount)
		return ErrInvalidAmount
	}

	// Parse transaction time
	transactionTime, err := time.Parse("2006-01-02 15:04:05", notification.TransactionTime)
	if err != nil {
		log.Printf("Failed to parse transaction time: %v", err)
		transactionTime = time.Now()
	}

	// Update payment with notification data
	transactionStatus := models.TransactionStatus(notification.TransactionStatus)
	fraudStatus := models.FraudStatus(notification.FraudStatus)

	payment.TransactionID = &notification.TransactionID
	payment.TransactionStatus = &transactionStatus
	payment.PaymentType = &notification.PaymentType
	payment.TransactionTime = &transactionTime
	payment.StatusMessage = &notification.StatusMessage
	payment.FraudStatus = &fraudStatus

	// Parse settlement time if present
	if notification.SettlementTime != nil && *notification.SettlementTime != "" {
		var settlementTime time.Time
		settlementTime, err = time.Parse("2006-01-02 15:04:05", *notification.SettlementTime)
		if err == nil {
			payment.SettlementTime = &settlementTime
		}
	}

	// Store full notification as metadata
	metadataBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification metadata: %v", err)
		payment.PaymentMetadata = datatypes.JSON("{}")
	} else {
		payment.PaymentMetadata = datatypes.JSON(metadataBytes)
	}

	// Update payment record
	if err := s.paymentRepo.Update(payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Update order status based on transaction status
	var newOrderStatus models.OrderStatus
	shouldUpdateOrder := true

	switch transactionStatus {
	case models.TransactionStatusSettlement:
		newOrderStatus = models.OrderStatusPreparing
		log.Printf("Payment settled for order: %s", notification.OrderID)
	case models.TransactionStatusPending:
		shouldUpdateOrder = false
		log.Printf("Payment pending for order: %s", notification.OrderID)
	case models.TransactionStatusExpire, models.TransactionStatusCancel, models.TransactionStatusDeny:
		newOrderStatus = models.OrderStatusCancelled
		log.Printf("Payment failed for order: %s, status: %s", notification.OrderID, transactionStatus)
	default:
		shouldUpdateOrder = false
		log.Printf("Unknown transaction status for order: %s, status: %s", notification.OrderID, transactionStatus)
	}

	// Update order status
	if shouldUpdateOrder && payment.Order != nil {
		if err := s.orderRepo.UpdateStatus(payment.OrderID, newOrderStatus); err != nil {
			log.Printf("Failed to update order status: %v", err)
			return fmt.Errorf("failed to update order status: %w", err)
		}
	}

	return nil
}

func (s *paymentService) VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	// Midtrans signature format: SHA512(order_id + status_code + gross_amount + server_key)
	input := orderID + statusCode + grossAmount + s.serverKey
	hash := sha512.Sum512([]byte(input))
	expectedSignature := hex.EncodeToString(hash[:])

	return expectedSignature == signatureKey
}
