package Model

type User struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"not null;size:30" json:"name"`
	Phone          string `gorm:"not null;unique" json:"phone"`
	Email          string `gorm:"not null;unique" json:"email"`
	Password       string `gorm:"not null" json:"password"`
	Province       string `json:"province"`
	City           string `json:"city"`
	Subdistrict    string `json:"subdistrict"`
	Address        string `gorm:"not null;size:100" json:"address"` // Detailed Address
	ProfilePicture string `json:"profile_picture"`
	WalletID       uint   `json:"wallet_id"`
	Wallet         Wallet `gorm:"ForeignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserRegister struct {
	Name        string `gorm:"not null;size:30" json:"name"`
	Phone       string `gorm:"not null;unique" json:"phone"`
	Email       string `gorm:"not null;unique" json:"email"`
	Password    string `gorm:"not null" json:"password"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Subdistrict string `json:"subdistrict"`
	Address     string `gorm:"not null;size:100" json:"address"` // Detailed Address
}

type UserUpdate struct {
	Name        string `gorm:"not null;size:30" json:"name"`
	Password    string `gorm:"not null" json:"password"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Subdistrict string `json:"subdistrict"`
	Address     string `gorm:"not null;size:100" json:"address"` // Detailed Address
}

type UserLogin struct {
	Phone    string `gorm:"not null;unique" json:"phone"`
	Password string `gorm:"not null" json:"password"`
}
