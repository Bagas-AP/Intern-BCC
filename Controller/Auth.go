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

		newUser := Model.User{
			Name:        input.Name,
			Phone:       input.Phone,
			Email:       input.Email,
			Password:    Hash(input.Password),
			Province:    input.Province,
			City:        input.City,
			Subdistrict: input.Subdistrict,
			Address:     input.Address,
			CreatedAt:   time.Now(),
		}

		if err := db.Create(&newUser).Error; err != nil {
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

			strToken, err := token.SignedString([]byte(os.Getenv("TOKEN")))
			if err != nil {
				c.JSON(http.StatusInternalServerError, Utils.FailedResponse(err.Error()))
				return
			}

			c.JSON(http.StatusOK, Utils.SucceededReponse("Success", gin.H{
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
