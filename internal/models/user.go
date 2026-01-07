package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleMember  UserRole = "member"
	RoleKiosk   UserRole = "kiosk"
	RoleBarista UserRole = "barista"
	RoleAdmin   UserRole = "admin"
)

type User struct {
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Phone     *string   `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	FullName  string    `gorm:"type:varchar(255);not null" json:"full_name"`
	Role      UserRole  `gorm:"type:varchar(20);not null" json:"role"`
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
}

func (User) TableName() string {
	return "users"
}
