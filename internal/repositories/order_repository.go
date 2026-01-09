package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound        = errors.New("order not found")
	ErrInvalidOrderStatus   = errors.New("invalid order status")
	ErrOrderNumberGenFailed = errors.New("failed to generate order number")
)

type OrderFilters struct {
	Status      *models.OrderStatus
	OrderSource *models.OrderSource
	StartDate   *time.Time
	EndDate     *time.Time
}

type OrderRepository interface {
	Create(order *models.Order, items []models.OrderItem) error
	FindByID(id uint) (*models.Order, error)
	FindByUUID(uuid uuid.UUID) (*models.Order, error)
	FindByOrderNumber(orderNumber string) (*models.Order, error)
	FindByUserID(userID uint, limit, offset int) ([]models.Order, int64, error)
	FindAll(filters OrderFilters, limit, offset int) ([]models.Order, int64, error)
	UpdateStatus(orderID uint, status models.OrderStatus) error

	GenerateOrderNumber() (string, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order, items []models.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Set OrderID for all items
		for i := range items {
			items[i].OrderID = order.ID
		}

		// Create all order items
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *orderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Payments").
		Where("id = ?", id).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUUID(uuid uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.
		Preload("Items").
		Preload("Items.Product").
		Where("uuid = ?", uuid).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Payments").
		Where("order_number = ?", orderNumber).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID uint, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// Count total orders for the user
	if err := r.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated orders
	err := r.db.
		Preload("Items").
		Preload("Items.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) FindAll(filters OrderFilters, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{})

	// Apply filters
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.OrderSource != nil {
		query = query.Where("order_source = ?", *filters.OrderSource)
	}

	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", *filters.StartDate)
	}

	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", *filters.EndDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated orders with preloads
	err := query.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Payments").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) UpdateStatus(orderID uint, status models.OrderStatus) error {
	updates := map[string]any{
		"status": status,
	}

	if status == models.OrderStatusCompleted {
		updates["completed_at"] = time.Now()
	}

	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Updates(updates).Error
}

func (r *orderRepository) GenerateOrderNumber() (string, error) {
	var orderNumber string

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Get current date in YYMMDD format
		today := time.Now().Format("060102")
		prefix := "MC-" + today + "-"

		// Get last order number for today
		var lastOrder models.Order
		err := tx.Where("order_number LIKE ?", prefix+"%").
			Order("order_number DESC").
			First(&lastOrder).Error

		sequence := 1
		if err == nil {
			// Extract sequence from last order number
			// Format: MC-YYMMDD-XXX
			parts := strings.Split(lastOrder.OrderNumber, "-")
			if len(parts) == 3 {
				seq, parseErr := strconv.Atoi(parts[2])
				if parseErr == nil {
					sequence = seq + 1
				}
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Generate order number with 3-digit sequence
		orderNumber = fmt.Sprintf("%s%03d", prefix, sequence)
		return nil
	})

	if err != nil {
		return "", ErrOrderNumberGenFailed
	}

	return orderNumber, nil
}
