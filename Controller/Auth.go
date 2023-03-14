package Controller

import (
	"bcc/Model"
	"bcc/Utils"
	"github.com/google/uuid"
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
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		newUser := Model.User{
			Name:        input.Name,
			Phone:       input.Phone,
			Email:       input.Email,
			Password:    Utils.Hash(input.Password),
			Province:    input.Province,
			City:        input.City,
			Subdistrict: input.Subdistrict,
			Address:     input.Address,
			CreatedAt:   time.Now(),
		}

		if err := db.Create(&newUser).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		var wallet Model.Wallet

		db.Where("user_id = ?", newUser.ID).First(&wallet)

		if wallet.ID != uuid.Nil {
			// The user already has a wallet
			return
		} else {
			// The user does not have a wallet
			newWallet := Model.Wallet{
				ID:     uuid.New(),
				UserID: newUser.ID,
			}

			if err := db.Create(&newWallet).Error; err != nil {
				Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
				return
			}
		}

		Utils.HttpRespSuccess(c, http.StatusCreated, "Account created", input)
	})
}

func Login(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	// user login
	r.POST("/login", func(c *gin.Context) {
		var input Model.UserLogin
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		login := Model.User{}
		if err := db.Where("phone = ?", input.Phone).First(&login).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		if login.Password == Utils.Hash(input.Password) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
				"id":  login.ID,
				"exp": time.Now().Add(time.Hour * 7 * 24).Unix(),
			})

			strToken, err := token.SignedString([]byte(os.Getenv("TOKEN")))
			if err != nil {
				Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusOK, "Parsed token", gin.H{
				"name":  login.Name,
				"token": strToken,
			})

		} else {
			Utils.HttpRespFailed(c, http.StatusForbidden, "Wrong password")
			return
		}
	})
}
