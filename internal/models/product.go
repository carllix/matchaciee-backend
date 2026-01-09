package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID              uint                   `gorm:"primaryKey;autoIncrement" json:"-"`
	UUID            uuid.UUID              `gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()" json:"id"`
	CategoryID      *uint                  `gorm:"index" json:"-"`
	Name            string                 `gorm:"type:varchar(255);not null" json:"name"`
	Slug            string                 `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description     *string                `gorm:"type:text" json:"description,omitempty"`
	BasePrice       float64                `gorm:"type:decimal(10,2);not null" json:"base_price"`
	PreparationTime int                    `gorm:"default:5" json:"preparation_time"`
	DisplayOrder    int                    `gorm:"default:0" json:"display_order"`
	IsAvailable     bool                   `gorm:"default:true" json:"is_available"`
	IsCustomizable  bool                   `gorm:"default:false" json:"is_customizable"`
	ImageURL        *string                `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	DeletedAt       *time.Time             `gorm:"index" json:"deleted_at,omitempty"`
	Category        *Category              `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
	Customizations  []ProductCustomization `gorm:"foreignKey:ProductID;references:ID" json:"customizations,omitempty"`
	CreatedAt       time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
