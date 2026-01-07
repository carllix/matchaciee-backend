package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OrderItem struct {
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Notes          *string        `gorm:"type:text" json:"notes,omitempty"`
	ProductID      *uint          `gorm:"index" json:"-"`
	Order          *Order         `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	Product        *Product       `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:SET NULL" json:"product,omitempty"`
	ProductName    string         `gorm:"type:varchar(255);not null" json:"product_name"`
	Customizations datatypes.JSON `gorm:"type:jsonb" json:"customizations,omitempty"`
	OrderID        uint           `gorm:"not null;index" json:"-"`
	Quantity       int            `gorm:"not null;default:1" json:"quantity"`
	UnitPrice      float64        `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal       float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"-"`
	UUID           uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()" json:"id"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
