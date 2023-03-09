package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func userCatering(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user")

	// get all catering
	r.GET("allCatering", Middleware.Authorization(), func(c *gin.Context) {
		var caterings []Model.Catering
		if err := db.Find(&caterings).Error; err != nil {
			c.JSON(http.StatusNotFound, Utils.FailedResponse(err.Error()))
			return
		}
	})

	// search catering by name
	r.POST("/searchLaundry", Middleware.Authorization(), func(c *gin.Context) {
		name, _ := c.GetQuery("name")

		var caterings []Model.Catering

		if res := db.Where("name LIKE ?", "%"+name+"%").Find(&caterings); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by name", gin.H{
			"data": caterings,
		}))
	})
}
