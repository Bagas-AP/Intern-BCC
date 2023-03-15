package Controller

import (
	"bcc/Model"
	"bcc/Utils"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
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
	r := q.Group("/api/user")
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

		forgotPW := Model.ResetPassword{
			UserID: user.ID,
			Phone:  user.Phone,
			Code:   generateRandomCode(),
		}

		var check Model.ResetPassword
		if err := db.Where("user_id = ?", user.ID).First(&check).Error; err == nil {
			Utils.HttpRespFailed(c, http.StatusForbidden, "You already requested a reset password")
			return
		}

		if err := db.Create(&forgotPW).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		if err := sendCodeToEmail(user.Email, forgotPW.Code); err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Reset password code sent", gin.H{
			"code": forgotPW.Code,
		})
	})
}

func sendCodeToEmail(email, code string) error {
	from := os.Getenv("EMAIL")              // Replace with your email address
	password := os.Getenv("EMAIL_PASSWORD") // Replace with your email password
	to := email

	// Message to be sent
	subject := "One-Time Code"
	body := fmt.Sprintf("Your one-time code is: %s", code)
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// SMTP configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}

func generateRandomCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 10
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}
