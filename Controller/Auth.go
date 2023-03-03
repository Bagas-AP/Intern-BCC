package Controller

import (
	"bcc/Model"
	"bcc/Utils"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	r.POST("/register", func(c *gin.Context) {
		var input Model.UserRegister
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusUnprocessableEntity, Utils.FailedResponse(err.Error()))
			return
		}

		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Utils.FailedResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Account created successfully", gin.H{
			"data": input,
		}))
	})
}

func Login(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	// user login
	r.POST("/login", func(c *gin.Context) {
		var input Model.UserLogin
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusUnprocessableEntity, Utils.FailedResponse(err.Error()))
			return
		}

		login := Model.User{}
		if err := db.Where("phone = ?", input.Phone).First(&login).Error; err != nil {
			c.JSON(http.StatusNotFound, Utils.FailedResponse(err.Error()))
			return
		}

		if login.Password == Hash(input.Password) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
				"id":  login.ID,
				"exp": time.Now().Add(time.Hour * 7 * 24).Unix(),
			})

			if err := godotenv.Load("../.env"); err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			strToken, err := token.SignedString([]byte(os.Getenv("TOKEN_G")))
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			c.JSON(http.StatusOK, Utils.SucceededReponse("Welcome, take your token", gin.H{
				"name":  login.Name,
				"token": strToken,
			}))
			return

		} else {
			c.JSON(http.StatusForbidden, Utils.FailedResponse("Wrong password"))
			return
		}
	})
}

func Hash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	pw := hex.EncodeToString(hash.Sum(nil))
	return pw
}
