package Model

type CateringTags struct {
	Catering   Catering `gorm:"ForeignKey:CateringID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CateringID uint     `json:"catering_id"`
	Tag        string   `gorm:"size:30; not null" json:"tag"`
}
