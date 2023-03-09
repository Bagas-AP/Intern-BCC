package Config

import (
	"bcc/Model"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MakeConnectionDatabase(data *Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s "+
		"password=%s "+
		"host=%s "+
		"TimeZone=Asia/Singapore "+
		"port=%s "+
		"dbname=%s", data.SupabaseUser, data.SupabasePassword, data.SupabaseHost, data.SupabasePort, data.SupabaseDbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&Model.User{},
		&Model.WalletTransaction{},
		&Model.Wallet{},
		&Model.Seller{},
		&Model.Laundry{},
		&Model.LaundryMenu{},
		&Model.LaundryPhotos{},
		&Model.Favorite{},
		&Model.UserOrder{},
	); err != nil {
		return nil, err
	}
	return db, nil
}
