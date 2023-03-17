package Controller

import (
	"bcc/Model"
	"bcc/Utils"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/google/uuid"

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

		log.Println("otw buat user")
		if err := db.Create(&newUser).Error; err != nil {
			log.Println("error di sini")
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		log.Println("user created but masuk ke wallet")
		var wallet Model.Wallet

		if err := db.Where("user_id = ?", newUser.ID).First(&wallet).Error; err == nil {
			log.Println("error di sini juga")
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

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

		var login Model.User
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

func ResetPassword(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	r.POST("/reset-password", func(c *gin.Context) {
		var input Model.ResetPasswordInput
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		var user Model.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		newRandomPassword := generateRandomPassword()

		user.Password = Utils.Hash(newRandomPassword)

		if err := db.Save(&user).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
		}

		if err := SendToEmail(db, input.Email, newRandomPassword); err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Password reset", nil)
	})
}

func SendToEmail(db *gorm.DB, email string, password string) error {
	var user Model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	auth := smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("EMAIL_PASS_APP"), "smtp.gmail.com")

	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Password Reset\r\n\r\n"+
			"Dear User,\r\n\r\n"+
			"We have received a request to reset your password for your account associated with this email. Please find below the details of the new password.\r\n\r\n"+
			"Phone: %s\r\n"+
			"Email: %s\r\n"+
			"New Password: %s\r\n\r\n"+
			"If you did not make this request, please contact our support team immediately.\r\n\r\n"+
			"Best regards,\r\n"+
			"%s",
		os.Getenv("EMAIL"), email, user.Phone, email, password, "TriServe",
	)

	err := smtp.SendMail("smtp.gmail.com"+":"+"587", auth, os.Getenv("EMAIL_FROM"), []string{email}, []byte(message))
	return err
}

func generateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 10
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}
