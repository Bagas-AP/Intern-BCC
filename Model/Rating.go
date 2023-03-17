package Model

type Rating struct {
	User      User    `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    uint    `json:"user_id"`
	Model     int     `json:"model"` // 1 -> laundry, 2 -> catering
	ServiceID int     `json:"service_id"`
	Rating    float64 `json:"rating"`
}

type InputRating struct {
	Rating float64 `json:"rating"`
}
