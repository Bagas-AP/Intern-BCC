package Model

type Favourite struct {
	User      User `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    uint `json:"user_id"`
	Model     int  `json:"model"` // 1 -> laundry, 2 -> catering
	ServiceID int  `json:"service_id"`
	MenuID    int  `json:"menu_id"`
}

type FavouriteResult struct {
	ModelID   int    `json:"model_id"`
	ServiceID int    `json:"service_id"`
	MenuID    int    `json:"menu_id"`
	MenuName  string `json:"menu_name"`
	MenuPrice int    `json:"menu_price"`
	MenuPhoto string `json:"menu_photo"`
}
