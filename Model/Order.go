package Model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	User      User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    uint      `json:"user_id"`
	Model     int       `json:"model"` // 1 -> laundry, 2 -> catering
	ServiceID int       `json:"service_id"`
	MenuID    int       `json:"menu_id"`
	// logicnya mulai di sini
	Status      int `json:"status"`
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
}
