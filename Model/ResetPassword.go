package Model

import "time"

type ResetPassword struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	User     User   `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Phone    string `gorm:"not null" json:"phone"`
	Code     string `gorm:"not null" json:"code"`
	ExpireAt time.Time
}

type ResetPasswordInput struct {
	Email string `gorm:"not null" json:"email"`
}

type ResetPasswordConfirmationInput struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `gorm:"not null" json:"code"`
	Password    string `gorm:"not null" json:"password"`
	ConfirmPass string `gorm:"not null" json:"confirm_pass"`
}
