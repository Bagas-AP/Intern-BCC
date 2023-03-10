package Model

type CateringMenu struct {
	ID          uint     `gorm:"primaryKey" json:"id"`
	Catering    Catering `gorm:"ForeignKey:CateringID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CateringID  uint     `json:"catering_id"`
	MenuIndex   int      `json:"menu_index"`
	Name        string   `gorm:"size:30;not null" json:"name"`
	Description string   `gorm:"size:255" json:"description"`
	Price       int      `json:"price"`
	Photo       string   `json:"photo"`
}

type CateringMenuDetailed struct {
	Catering       Catering     `gorm:"ForeignKey:CateringID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CateringID     uint         `json:"catering_id"`
	CateringMenu   CateringMenu `gorm:"ForeignKey:CateringMenuID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CateringMenuID uint         `json:"catering_menu_id"`
	MenuMonday     string       `gorm:"not null" json:"menu_monday"`
	MenuTuesday    string       `gorm:"not null" json:"menu_tuesday"`
	MenuWednesday  string       `gorm:"not null" json:"menu_wednesday"`
	MenuThursday   string       `gorm:"not null" json:"menu_thursday"`
	MenuFriday     string       `gorm:"not null" json:"menu_friday"`
	MenuSaturday   string       `gorm:"not null" json:"menu_saturday"`
	MenuSunday     string       `gorm:"not null" json:"menu_sunday"`
}
