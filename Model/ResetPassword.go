package Model

type ResetPassword struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	User   User   `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint   `json:"user_id"`
	Phone  string `gorm:"not null;unique" json:"phone"`
	Code   string `gorm:"not null" json:"code"`
}

type ResetPasswordInput struct {
	Email string `gorm:"not null" json:"email"`
}
