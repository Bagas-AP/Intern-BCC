package Model

type LaundryPhotos struct {
	Laundry   Laundry `gorm:"ForeignKey:LaundryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LaundryID uint    `json:"laundry_id"`
	Link      string  `json:"link"`
}
