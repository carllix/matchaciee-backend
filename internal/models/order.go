package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderSource string

const (
	OrderSourceGuest  OrderSource = "guest"
	OrderSourceMember OrderSource = "member"
	OrderSourceKiosk  OrderSource = "kiosk"
)

type Order struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"-"`
	UUID         uuid.UUID   `gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()" json:"id"`
	OrderNumber  string      `gorm:"type:varchar(20);uniqueIndex;not null" json:"order_number"`
	UserID       *uint       `gorm:"index" json:"-"`
	CustomerName string      `gorm:"type:varchar(255);not null" json:"customer_name"`
	Status       OrderStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	OrderSource  OrderSource `gorm:"type:varchar(20);not null;index" json:"order_source"`
	Subtotal     float64     `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Tax          float64     `gorm:"type:decimal(10,2);default:0" json:"tax"`
	Total        float64     `gorm:"type:decimal(10,2);not null" json:"total"`
	QueueNumber  *int        `gorm:"type:int" json:"queue_number,omitempty"`
	Notes        *string     `gorm:"type:text" json:"notes,omitempty"`
	CompletedAt  *time.Time  `json:"completed_at,omitempty"`
	User         *User       `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	Items        []OrderItem `gorm:"foreignKey:OrderID;references:ID" json:"items,omitempty"`
	Payments     []Payment   `gorm:"foreignKey:OrderID;references:ID" json:"payments,omitempty"`
	CreatedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP;index" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Order) TableName() string {
	return "orders"
}
