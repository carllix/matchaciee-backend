package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	UpdatedAt       time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedAt       time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ImageURL        *string                `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	Category        *Category              `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
	CategoryID      *uint                  `gorm:"index" json:"-"`
	Description     *string                `gorm:"type:text" json:"description,omitempty"`
	DeletedAt       *time.Time             `gorm:"index" json:"deleted_at,omitempty"`
	Slug            string                 `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Name            string                 `gorm:"type:varchar(255);not null" json:"name"`
	Customizations  []ProductCustomization `gorm:"foreignKey:ProductID;references:ID" json:"customizations,omitempty"`
	PreparationTime int                    `gorm:"default:5" json:"preparation_time"`
	DisplayOrder    int                    `gorm:"default:0" json:"display_order"`
	BasePrice       float64                `gorm:"type:decimal(10,2);not null" json:"base_price"`
	ID              uint                   `gorm:"primaryKey;autoIncrement" json:"-"`
	UUID            uuid.UUID              `gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()" json:"id"`
	IsAvailable     bool                   `gorm:"default:true" json:"is_available"`
	IsCustomizable  bool                   `gorm:"default:false" json:"is_customizable"`
}

func (Product) TableName() string {
	return "products"
}
