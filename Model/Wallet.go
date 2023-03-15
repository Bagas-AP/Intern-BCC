package Model

import (
	"github.com/google/uuid"
	"time"
)

type Wallet struct {
	ID               uuid.UUID          `gorm:"primaryKey" json:"id"`
	User             User               `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID           uint               `json:"user_id"`
	WalletCategories []WalletCategories `gorm:"ForeignKey:WalletID"`
}

type WalletCategories struct {
	ID                 uuid.UUID           `gorm:"primaryKey" json:"id"`
	WalletID           uuid.UUID           `gorm:"type:uuid;index" json:"wallet_id"`
	Wallet             Wallet              `gorm:"foreignKey:WalletID"`
	User               User                `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID             uint                `json:"user_id"`
	Index              int                 `json:"index"`
	Name               string              `json:"name"`
	Budget             int                 `json:"budget"`
	Balance            int                 `json:"balance"`
	WalletTransactions []WalletTransaction `gorm:"foreignKey:WalletCategoryID"`
}

type WalletTransaction struct {
	ID               uuid.UUID        `gorm:"primaryKey" json:"id"`
	User             User             `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	WalletID         uuid.UUID        `gorm:"type:uuid;index" json:"wallet_id"`
	Wallet           Wallet           `gorm:"foreignKey:WalletID"`
	UserID           uint             `json:"user_id"`
	WalletCategoryID uuid.UUID        `gorm:"type:uuid;index" json:"wallet_category_id"`
	WalletCategory   WalletCategories `gorm:"foreignKey:WalletCategoryID"`
	Expense          string           `json:"expense"`
	Amount           int              `json:"amount"`
	CreatedAt        time.Time
}

type WalletCategoryInput struct {
	Name   string `json:"name"`
	Budget int    `json:"budget"`
}

type WalletTransactionInput struct {
	Expense string `json:"expense"`
	Amount  int    `json:"amount"`
}
