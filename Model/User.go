package Model

import "time"

type User struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"not null;size:30;binding:required" json:"name"`
	Phone          string `gorm:"not null;unique;binding:required" json:"phone"`
	Email          string `gorm:"not null;unique;binding:required" json:"email"`
	Password       string `gorm:"not null;binding:required" json:"password"`
	Province       string `gorm:"not null;binding:required" json:"province"`
	City           string `gorm:"not null;binding:required" json:"city"`
	Subdistrict    string `gorm:"not null;binding:required" json:"subdistrict"`
	Address        string `gorm:"not null;size:100;binding:required" json:"address"`
	ProfilePicture string `json:"profile_picture"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserRegister struct {
	Name        string `gorm:"not null;size:30" json:"name"`
	Phone       string `gorm:"not null;unique" json:"phone"`
	Email       string `gorm:"not null;unique" json:"email"`
	Password    string `gorm:"not null" json:"password"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Subdistrict string `json:"subdistrict"`
	Address     string `gorm:"not null;size:100" json:"address"`
}

type UserUpdate struct {
	Name           string `gorm:"not null;size:30" json:"name"`
	Phone          string `gorm:"not null;unique" json:"phone"`
	Email          string `gorm:"not null;unique" json:"email"`
	Password       string `gorm:"not null" json:"password"`
	Province       string `json:"province"`
	City           string `json:"city"`
	Subdistrict    string `json:"subdistrict"`
	Address        string `gorm:"not null;size:100" json:"address"`
	ProfilePicture string `json:"profile_picture"`
}

type UserLogin struct {
	Phone    string `gorm:"not null;unique" json:"phone"`
	Password string `gorm:"not null" json:"password"`
}
