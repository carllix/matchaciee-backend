package models

import (
	"time"
)

type RefreshToken struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"-"`
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	Token           string     `gorm:"type:varchar(500);uniqueIndex;not null" json:"-"`
	ExpiresAt       time.Time  `gorm:"not null;index" json:"expires_at"`
	RevokedAt       *time.Time `gorm:"index" json:"revoked_at,omitempty"`
	ReplacedByToken *string    `gorm:"type:varchar(500)" json:"-"`
	User            *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	CreatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (rt *RefreshToken) IsValid() bool {
	return rt.RevokedAt == nil && rt.ExpiresAt.After(time.Now())
}
