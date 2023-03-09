package Model

type CateringPhoto struct {
	Catering   Catering `gorm:"ForeignKey:CateringID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CateringID uint     `json:"catering_id"`
	Link       string   `json:"link"`
}
