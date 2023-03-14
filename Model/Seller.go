package Model

import "time"

type Seller struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"size:30;not null;binding:required" json:"name"`
	Password  string `gorm:"size:30;not null;binding:required" json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
