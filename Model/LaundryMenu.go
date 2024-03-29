package Model

type LaundryMenu struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Laundry     Laundry `gorm:"ForeignKey:LaundryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LaundryID   uint    `json:"laundry_id"`
	MenuIndex   int     `json:"menu_index"`
	Name        string  `gorm:"size:1000000;not null" json:"name"`
	Description string  `gorm:"size:255" json:"description"`
	Price       int     `json:"price"`
	Photo       string  `json:"photo"`
}
