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
	UpdatedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP;index" json:"created_at"`
	QueueNumber  *int        `gorm:"type:int" json:"queue_number,omitempty"`
	User         *User       `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	UserID       *uuid.UUID  `gorm:"type:uuid;index" json:"user_id,omitempty"`
	CompletedAt  *time.Time  `json:"completed_at,omitempty"`
	Notes        *string     `gorm:"type:text" json:"notes,omitempty"`
	OrderSource  OrderSource `gorm:"type:varchar(20);not null;index" json:"order_source"`
	Status       OrderStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	CustomerName string      `gorm:"type:varchar(255);not null" json:"customer_name"`
	OrderNumber  string      `gorm:"type:varchar(20);uniqueIndex;not null" json:"order_number"`
	Items        []OrderItem `gorm:"foreignKey:OrderID;references:ID" json:"items,omitempty"`
	Payments     []Payment   `gorm:"foreignKey:OrderID;references:ID" json:"payments,omitempty"`
	Total        float64     `gorm:"type:decimal(10,2);not null" json:"total"`
	Tax          float64     `gorm:"type:decimal(10,2);default:0" json:"tax"`
	Subtotal     float64     `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
}

func (Order) TableName() string {
	return "orders"
}
