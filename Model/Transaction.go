package Model

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID uuid.UUID `gorm:"primaryKey" json:"id"`
	// di order itu cuma userid, transaction id, pokoknya cuma kek jembatan anatar aja
	User          User   `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID        uint   `json:"user_id"`
	Model         int    `json:"model"` // 1 -> laundry, 2 -> catering
	ServiceID     int    `json:"service_id"`
	MenuID        int    `json:"menu_id"`
	Quantity      int    `json:"quantity"`
	Total         int    `json:"total"`
	PaymentMethod string `json:"payment_method"`
	PaymentPhoto  string `json:"payment_photo"`
	CreatedAt     time.Time
	UpdateAt      time.Time
}

type TransactionInput struct {
	Quantity      int    `json:"quantity"`
	PaymentMethod string `json:"payment_method"`
}

type TransactionResult struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	MenuName        string `json:"menu_name"`
	MenuPrice       string `json:"menu_price"`
	MenuExtendedFee int    `json:"menu_extended_fee"`
	Total           int    `json:"total"`
}
