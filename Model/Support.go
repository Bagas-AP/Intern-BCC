package Model

import "time"

type Support struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	User      User   `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    uint   `json:"user_id"`
	Help      string `json:"help"`
	CreatedAt time.Time
}

type SupportInput struct {
	Help string `json:"help"`
}
