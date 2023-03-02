package Controller

import (
	"bcc/Model"
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
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Error when binding JSON",
				"error":   err.Error(),
			})
			return
		}

		if err := db.Create(&input); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong with student creation",
				"error":   err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Account created successfully",
			"error":   nil,
			"data": gin.H{
				"nama": input.Name,
			},
		})
	})
}

func Login(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	// user login
	r.POST("/login", func(c *gin.Context) {
		type body struct {
			Phone    string `json:"phone"`
			Password string `json:"password"`
		}
		var input body
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Error when binding JSON",
				"error":   err.Error(),
			})
			return
		}
		login := Model.User{}
		if err := db.Where("phone = ?", input.Phone).First(&login).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Username does not exist",
				"error":   err.Error(),
			})
			return
		}
		if login.Password == Hash(input.Password) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
				"id":  login.ID,
				"exp": time.Now().Add(time.Hour * 7 * 24).Unix(),
			})
			godotenv.Load("../.env")
			strToken, err := token.SignedString([]byte(os.Getenv("TOKEN_G")))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when loading token",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Welcome, take your token",
				"data": gin.H{
					"name":  login.Name,
					"token": strToken,
				},
			})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Wrong password",
			})
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
