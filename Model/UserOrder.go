package Model

import "github.com/google/uuid"

type UserOrder struct {
	ID              uuid.UUID   `gorm:"primaryKey" json:"id"`
	User            User        `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID          uint        `json:"user_id"`
	Laundry         Laundry     `gorm:"ForeignKey:LaundryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LaundryID       uint        `json:"laundry_id"`
	LaundryMenu     LaundryMenu `gorm:"ForeignKey:MenuID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MenuID          uint        `json:"menu_id"`
	LaundryQuantity int         `json:"laundry_quantity"`
}

type InputLaundryOrder struct {
	LaundryID     int `json:"laundry_id"`
	LaundryMenuID int `json:"laundry_menu_id"`
	Quantity      int `json:"quantity"`
}
