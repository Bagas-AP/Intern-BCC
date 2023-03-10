package Middleware

import (
	"bcc/Utils"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		header = header[len("Bearer "):]

		token, err := jwt.Parse(header, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN")), nil
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(err.Error()))
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("id", uint(claims["id"].(float64)))
			c.Next()
			return
		} else {
			log.Println("masuk ke else")
			c.JSON(http.StatusForbidden, Utils.FailedResponse(err.Error()))
			c.Abort()
			return
		}
	}
}
