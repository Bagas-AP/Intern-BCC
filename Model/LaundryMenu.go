package Model

type LaundryMenu struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Laundry     Laundry `gorm:"ForeignKey:LaundryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LaundryID   uint    `json:"laundry_id"`
	Name        string  `gorm:"size:30;not null" json:"name"`
	Description string  `gorm:"size:255" json:"description"`
	Price       int     `json:"price"`
	Photo       string  `json:"photo"`
}
