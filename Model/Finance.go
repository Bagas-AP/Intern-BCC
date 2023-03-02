package Model

import "time"

type Wallet struct {
	ID    uint `gorm:"primary_key" json:"id"`
	Total int  `json:"total"`
}

type WalletTransaction struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	User      User   `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    uint   `json:"user_id"`
	Wallet    Wallet `gorm:"ForeignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	WalletID  uint   `json:"wallet_id"`
	Expense   string `gorm:"size:100" json:"expense"`
	Amount    int    `json:"amount"`
	CreatedAt time.Time
}
