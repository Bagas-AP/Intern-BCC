package Model

type Catering struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	Seller     Seller  `gorm:"ForeignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SellerID   uint    `json:"seller_id"`
	Name       string  `gorm:"size:255;not null" json:"name"`
	Address    string  `gorm:"size:255;not null" json:"address"`
	Rating     float64 `json:"rating"`
	Instagram  string  `json:"instagram"`
	Phone      string  `json:"phone"`
	PriceRange string  `json:"price_range"`
	Picture    string  `json:"picture"`
}

type CateringSearch struct {
	Name string `json:"name"`
}
