package Model

type LaundryTags struct {
	Laundry   Laundry `gorm:"ForeignKey:LaundryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LaundryID uint    `json:"laundry_id"`
	Tag       string  `gorm:"size:30" json:"tag"`
}
