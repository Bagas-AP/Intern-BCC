package App

import (
	"bcc/Config"
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

	databaseConf, err := Config.NewDatabase()
	if err != nil {
		panic(err.Error())
	}
	db, err := Config.MakeConnectionDatabase(databaseConf)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Database Connected")

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	fmt.Println(db)
	fmt.Println("Cuma print db kok")

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
