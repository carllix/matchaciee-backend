package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Description  *string   `gorm:"type:text" json:"description,omitempty"`
	ImageURL     *string   `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Slug         string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"-"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	UUID         uuid.UUID `gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()" json:"id"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
}

func (Category) TableName() string {
	return "categories"
}
