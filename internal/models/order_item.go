package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OrderItem struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"-"`
	UUID           uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()" json:"id"`
	OrderID        uint           `gorm:"not null;index" json:"-"`
	ProductID      *uint          `gorm:"index" json:"-"`
	ProductName    string         `gorm:"type:varchar(255);not null" json:"product_name"`
	Quantity       int            `gorm:"not null;default:1" json:"quantity"`
	UnitPrice      float64        `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal       float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Customizations datatypes.JSON `gorm:"type:jsonb" json:"customizations,omitempty"`
	Notes          *string        `gorm:"type:text" json:"notes,omitempty"`
	Order          *Order         `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	Product        *Product       `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:SET NULL" json:"product,omitempty"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
