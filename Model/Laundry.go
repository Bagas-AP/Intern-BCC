package Model

type Laundry struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Seller      Seller  `gorm:"ForeignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SellerID    uint    `json:"seller_id"`
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"size:255" json:"description"`
	Address     string  `gorm:"size:255;not null" json:"address"`
	Rating      float64 `json:"rating"`
	Instagram   string  `json:"instagram"`
	Phone       string  `json:"phone"`
	PriceRange  string  `json:"price_range"`
}

type LaundrySearch struct {
	Name string `json:"name"`
}
