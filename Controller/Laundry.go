package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserLaundry(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user")
	// get all laundry
	r.GET("/allLaundry", Middleware.Authorization(), func(c *gin.Context) {
		var laundries []Model.Laundry
		if err := db.Find(&laundries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Utils.FailedResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get all laundry", gin.H{
			"data": laundries,
		}))
	})

	// seach laundry by name
	r.POST("/searchLaundry", Middleware.Authorization(), func(c *gin.Context) {
		name, _ := c.GetQuery("name")

		var laundries []Model.Laundry

		if res := db.Where("name LIKE ?", "%"+name+"%").Find(&laundries); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by name", gin.H{
			"data": laundries,
		}))
	})

	// get laundry by id
	r.GET("/laundry/:id", Middleware.Authorization(), func(c *gin.Context) {
		id := c.Param("id")

		var laundry Model.Laundry
		if res := db.Where("id = ?", id).First(&laundry); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		var menus []Model.LaundryMenu
		if res := db.Where("laundry_id = ?", laundry.ID).Find(&menus); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by id", gin.H{
			"laundry_name":       laundry.Name,
			"laundry_address":    laundry.Address,
			"laundry_phone":      laundry.Phone,
			"laundry_instagram":  laundry.Instagram,
			"laundry_rating":     laundry.Rating,
			"laundry_priceRange": laundry.PriceRange,
			"menus":              menus,
		}))
	})
}
