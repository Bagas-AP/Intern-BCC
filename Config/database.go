package Config

import (
	"bcc/Model"
	"fmt"
	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MakeSupaBaseClient() supabasestorageuploader.SupabaseClientService {
	SupaBaseClient := supabasestorageuploader.NewSupabaseClient(
		"https://llghldcosbvakddztrpt.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImxsZ2hsZGNvc2J2YWtkZHp0cnB0Iiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY3NzUwMjEzOCwiZXhwIjoxOTkzMDc4MTM4fQ.1fHxjZfPYIaTEEzEh9aUg5_T0yh-cyq5KG9KAbxt8C4",
		"picture",
		"",
	)
	return SupaBaseClient
}

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
		//&Model.WalletTransaction{},
		//&Model.Wallet{},
		&Model.Seller{},
		&Model.Laundry{},
		&Model.LaundryMenu{},
		&Model.LaundryPhotos{},
		&Model.Favourite{},
		&Model.Order{},
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
		&Model.CateringTags{},
		&Model.Favourite{},
		&Model.Transaction{},
		&Model.Order{},
		&Model.Rating{},
		&Model.Wallet{},
		&Model.WalletCategories{},
		&Model.WalletTransaction{},
		&Model.ResetPassword{},
	); err != nil {
		return nil, err
	}
	log.Println("db migrated")
	return db, nil
}
