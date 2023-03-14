package Model

import (
	"github.com/google/uuid"
	"time"
)

type Wallet struct {
	ID     uuid.UUID `gorm:"primaryKey" json:"id"`
	User   User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint      `json:"user_id"`
}

type WalletCategories struct {
	ID       uuid.UUID `gorm:"primaryKey" json:"id"`
	Wallet   Wallet    `gorm:"ForeignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	WalletID uuid.UUID `json:"wallet_id"`
	Index    int       `json:"index"`
	Name     string    `json:"name"`
	Budget   int       `json:"budget"`
	Balance  int       `json:"balance"`
}

type WalletTransaction struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	User      User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    uint      `json:"user_id"`
	Wallet    Wallet    `gorm:"ForeignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	WalletID  uuid.UUID `json:"wallet_id"`
	Expense   string    `json:"expense"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time
}

//type WalletTransaction struct {
//	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
//	User      User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//	UserID    uint      `json:"user_id"`
//	Wallet    Wallet    `gorm:"ForeignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//	WalletID  uint      `json:"wallet_id"`
//	Expense   string    `gorm:"size:100" json:"expense"`
//	Amount    int       `json:"amount"`
//	CreatedAt time.Time
//}
