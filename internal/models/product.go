package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	CreatedAt       time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Description     *string                `gorm:"type:text" json:"description,omitempty"`
	DeletedAt       *time.Time             `gorm:"index" json:"deleted_at,omitempty"`
	Category        *Category              `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
	CategoryID      *uuid.UUID             `gorm:"type:uuid;index" json:"category_id,omitempty"`
	ImageURL        *string                `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	Name            string                 `gorm:"type:varchar(255);not null" json:"name"`
	Slug            string                 `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Customizations  []ProductCustomization `gorm:"foreignKey:ProductID;references:ID" json:"customizations,omitempty"`
	PreparationTime int                    `gorm:"default:5" json:"preparation_time"`
	DisplayOrder    int                    `gorm:"default:0" json:"display_order"`
	BasePrice       float64                `gorm:"type:decimal(10,2);not null" json:"base_price"`
	ID              uuid.UUID              `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	IsCustomizable  bool                   `gorm:"default:false" json:"is_customizable"`
	IsAvailable     bool                   `gorm:"default:true" json:"is_available"`
}

func (Product) TableName() string {
	return "products"
}
