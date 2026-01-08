package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductCustomization struct {
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Product           *Product  `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE" json:"product,omitempty"`
	CustomizationType string    `gorm:"type:varchar(50);not null" json:"customization_type"`
	OptionName        string    `gorm:"type:varchar(100);not null" json:"option_name"`
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"-"`
	ProductID         uint      `gorm:"not null;index" json:"-"`
	PriceModifier     float64   `gorm:"type:decimal(10,2);default:0" json:"price_modifier"`
	DisplayOrder      int       `gorm:"default:0" json:"display_order"`
	UUID              uuid.UUID `gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()" json:"id"`
}

func (ProductCustomization) TableName() string {
	return "product_customizations"
}
