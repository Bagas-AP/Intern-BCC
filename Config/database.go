package Config

import (
	"bcc/Model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func MakeSupaBaseConnectionDatabase(data *Database) (*gorm.DB, error) {
	// using supabase
	dsn := fmt.Sprintf("user=%s "+
		"password=%s "+
		"host=%s "+
		"TimeZone=Asia/Singapore "+
		"port=%s "+
		"dbname=%s",
		data.SupabaseUser, data.SupabasePassword, data.SupabaseHost, data.SupabasePort, data.SupabaseDbName)
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

func MakeLocalhostConnectionDatabase(data *DBLocal) (*gorm.DB, error) {
	// using localhost
	db, err := gorm.Open(
		mysql.Open(
			fmt.Sprintf(
				"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
				data.DbUser, data.DbPassword, data.DbHost, data.DbName,
			),
		),
		&gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := db.AutoMigrate(
		&Model.User{},
		&Model.Seller{},
		&Model.Laundry{},
		&Model.LaundryMenu{},
		&Model.LaundryTags{},
		&Model.LaundryPhotos{},
		&Model.Catering{},
		&Model.CateringMenu{},
		&Model.CateringMenuDetailed{},
		&Model.Favorite{},
		&Model.UserOrder{},
		//&Model.WalletTransaction{},
		//&Model.Wallet{},
	); err != nil {
		return nil, err
	}
	return db, nil
}
