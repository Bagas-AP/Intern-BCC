package Model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	User        User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID      uint      `json:"user_id"`
	Model       int       `json:"model"` // 1 -> laundry, 2 -> catering
	ServiceID   int       `json:"service_id"`
	MenuID      int       `json:"menu_id"`
	Status      int       `json:"status"`
	CreatedAt   time.Time
	UpdatedAt   time.Time // -> payed at
	CompletedAt time.Time
}

type OrderResult struct {
	ModelID     int    `json:"model_id"`
	ServiceID   int    `json:"service_id"`
	MenuID      int    `json:"menu_id"`
	ServiceName string `json:"service_name"`
	MenuName    string `json:"menu_name"`
	MenuPhoto   string `json:"menu_photo"`
	Status      int    `json:"status"`
	SellerPhone string `json:"seller_phone"`
}

type OrderInputConfirm struct {
	Confirm int `json:"confirm"`
}
