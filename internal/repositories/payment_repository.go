package repositories

import (
	"errors"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrPaymentDuplicate = errors.New("payment with this midtrans order ID already exists")
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByUUID(uuid uuid.UUID) (*models.Payment, error)
	FindByMidtransOrderID(midtransOrderID string) (*models.Payment, error)
	FindByOrderID(orderID uint) ([]models.Payment, error)
	Update(payment *models.Payment) error
	UpdateTransactionStatus(paymentID uint, status models.TransactionStatus) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByUUID(uuid uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.
		Preload("Order").
		Where("uuid = ?", uuid).
		First(&payment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByMidtransOrderID(midtransOrderID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.
		Preload("Order").
		Where("midtrans_order_id = ?", midtransOrderID).
		First(&payment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByOrderID(orderID uint) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&payments).Error

	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) UpdateTransactionStatus(paymentID uint, status models.TransactionStatus) error {
	return r.db.Model(&models.Payment{}).
		Where("id = ?", paymentID).
		Update("transaction_status", status).Error
}
