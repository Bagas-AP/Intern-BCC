package App

import (
	"bcc/Config"
	"bcc/Controller"
	"bcc/Middleware"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	// supabase
	//databaseConf, err := Config.NewDatabase()
	//if err != nil {
	//	panic(err.Error())
	//}
	//db, err := Config.MakeSupaBaseConnectionDatabase(databaseConf)
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println("Database Connected")

	// localhost
	databaseConf, err := Config.NewDBLocal()
	if err != nil {
		panic(err.Error())
	}
	db, err := Config.MakeLocalhostConnectionDatabase(databaseConf)
	fmt.Println("Database Connected")
	r := gin.Default()

	// cors
	r.Use(Middleware.CORS())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	Controller.Register(db, r)
	Controller.Login(db, r)
	Controller.UserProfile(db, r)
	Controller.UserLaundry(db, r)
	Controller.UserCatering(db, r)
	Controller.UserFavourite(db, r)
	Controller.UserOrder(db, r)
	Controller.UserTransaction(db, r)

	if err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		panic(err.Error())
	}
}
