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
	Controller.Profile(db, r)
	Controller.UserLaundry(db, r)
	Controller.UserCatering(db, r)

	fmt.Println(db)
	fmt.Println("Cuma print db kok")

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
