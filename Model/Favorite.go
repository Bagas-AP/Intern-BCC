package Model

type Favourite struct {
	User      User `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    uint `json:"user_id"`
	Model     int  `json:"model"` // 1 -> laundry, 2 -> catering
	ServiceID int  `json:"service_id"`
	MenuID    int  `json:"menu_id"`
}

type FavouriteResult struct {
	//ServiceID uint   `json:"service_id"`
	//MenuID    uint   `json:"menu_id"`
	MenuName  string `json:"menu_name"`
	MenuPrice int    `json:"menu_price"`
	MenuPhoto string `json:"menu_photo"`
}
