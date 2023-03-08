package Model

type Catering struct {
	// ID          uint   `gorm:"primaryKey" json:"id"`
	// Seller      Seller `gorm:"ForeignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// SellerID    uint   `json:"seller_id"`
	// Name        string `gorm:"size:30;not null" json:"name"`
	// Description string `gorm:"size:255" json:"description"`
	// Location    string `gorm:"size:255;not null" json:"location"`
	// Rating      int    `json:"rating"`
	// Instagram   string `json:"instagram"`
	// Phone       string `json:"phone"`
	// PriceRange  string `json:"price_range"`
}
