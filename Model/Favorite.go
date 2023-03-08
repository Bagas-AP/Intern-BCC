package Model

type Favorite struct {
	User          User        `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID        uint        `json:"user_id"`
	Laundry       Laundry     `gorm:"ForeignKey:LaundryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LaundryID     uint        `json:"laundry_id"`
	LaundryMenu   LaundryMenu `gorm:"ForeignKey:LaundryMenuID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LaundryMenuID uint        `json:"laundry_menu_id"`
	// Catering   Catering `gorm:"ForeignKey:CateringID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// CateringID uint     `json:"catering_id"`
	// CateringMenu   CateringMenu `gorm:"ForeignKey:CateringMenuID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// CateringMenuID uint        `json:"catering_menu_id"`
}
