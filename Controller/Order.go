package Controller

import (
	"bcc/Utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func UserOrder(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/order")
	r.GET("/all", func(c *gin.Context) {
		Utils.HttpRespSuccess(c, http.StatusOK, "Alive", nil)
	})
}
