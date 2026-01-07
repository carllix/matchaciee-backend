package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusSettlement TransactionStatus = "settlement"
	TransactionStatusExpire     TransactionStatus = "expire"
	TransactionStatusCancel     TransactionStatus = "cancel"
	TransactionStatusDeny       TransactionStatus = "deny"
	TransactionStatusRefund     TransactionStatus = "refund"
)

type FraudStatus string

const (
	FraudStatusAccept    FraudStatus = "accept"
	FraudStatusChallenge FraudStatus = "challenge"
	FraudStatusDeny      FraudStatus = "deny"
)

type Payment struct {
	CreatedAt         time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	TransactionID     *string            `gorm:"type:varchar(100)" json:"transaction_id,omitempty"`
	PaymentType       *string            `gorm:"type:varchar(50)" json:"payment_type,omitempty"`
	TransactionStatus *TransactionStatus `gorm:"type:varchar(50);index" json:"transaction_status,omitempty"`
	TransactionTime   *time.Time         `json:"transaction_time,omitempty"`
	SettlementTime    *time.Time         `json:"settlement_time,omitempty"`
	FraudStatus       *FraudStatus       `gorm:"type:varchar(50)" json:"fraud_status,omitempty"`
	StatusMessage     *string            `gorm:"type:text" json:"status_message,omitempty"`
	Order             *Order             `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	MidtransOrderID   string             `gorm:"type:varchar(100);uniqueIndex;not null" json:"midtrans_order_id"`
	PaymentMetadata   datatypes.JSON     `gorm:"type:jsonb" json:"payment_metadata,omitempty"`
	GrossAmount       float64            `gorm:"type:decimal(10,2);not null" json:"gross_amount"`
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID           uuid.UUID          `gorm:"type:uuid;not null;index" json:"order_id"`
}

func (Payment) TableName() string {
	return "payments"
}
