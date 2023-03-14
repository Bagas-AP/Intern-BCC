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
	UserAddress   string `json:"user_address"`
	Model         int    `json:"model"` // 1 -> laundry, 2 -> catering
	ServiceID     int    `json:"service_id"`
	MenuID        int    `json:"menu_id"`
	Notes         string `json:"notes"`
	Quantity      int    `json:"quantity"`
	Total         int    `json:"total"`
	PaymentMethod int    `json:"payment_method"`
	PaymentProof  string `json:"payment_proof"`
	CreatedAt     time.Time
	UpdateAt      time.Time
}

type TransactionPayment struct {
	PaymentDeadline time.Time `json:"payment_deadline"`
	AccountNumber   string    `json:"account_number"`
	AccountName     string    `json:"account_name"`
	Total           int       `json:"total"`
}

type TransactionChangeAddress struct {
	Province    string `gorm:"not null;binding:required" json:"province"`
	City        string `gorm:"not null;binding:required" json:"city"`
	Subdistrict string `gorm:"not null;binding:required" json:"subdistrict"`
	Address     string `gorm:"not null;size:100;binding:required" json:"address"` // Detailed Address
}

type TransactionInputNotes struct {
	Notes string `json:"notes"`
}

type TransactionInputQuantity struct {
	Quantity int `gorm:"binding:required" json:"quantity"`
}

type TransactionInputPayment struct {
	PaymentMethod int `gorm:"binding:required" json:"payment_method"`
}

type TransactionResult struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	MenuName        string `json:"menu_name"`
	MenuPrice       int    `json:"menu_price"`
	MenuExtendedFee int    `json:"menu_extended_fee"`
	Total           int    `json:"total"`
}
