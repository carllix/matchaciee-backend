package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"-"`
	UUID         uuid.UUID `gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Slug         string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description  *string   `gorm:"type:text" json:"description,omitempty"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	ImageURL     *string   `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
