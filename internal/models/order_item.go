package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OrderItem struct {
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ProductID      *uuid.UUID     `gorm:"type:uuid;index" json:"product_id,omitempty"`
	Notes          *string        `gorm:"type:text" json:"notes,omitempty"`
	Order          *Order         `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	Product        *Product       `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:SET NULL" json:"product,omitempty"`
	ProductName    string         `gorm:"type:varchar(255);not null" json:"product_name"`
	Customizations datatypes.JSON `gorm:"type:jsonb" json:"customizations,omitempty"`
	Quantity       int            `gorm:"not null;default:1" json:"quantity"`
	UnitPrice      float64        `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal       float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"order_id"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
