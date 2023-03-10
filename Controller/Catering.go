package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserCatering(db *gorm.DB, q *gin.Engine) {
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
	r.POST("/searchCatering", Middleware.Authorization(), func(c *gin.Context) {
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

	// get catering by id and all its menu
	r.GET("/catering/:catering_id", Middleware.Authorization(), func(c *gin.Context) {
		id := c.Param("catering_id")

		var catering Model.Catering
		if res := db.Where("id = ?", id).First(&catering); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		var menus []Model.CateringMenu
		if res := db.Where("catering_id = ?", catering.ID).Find(&menus); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by id", gin.H{
			"laundry_name":        catering.Name,
			"catering_address":    catering.Address,
			"catering_phone":      catering.Phone,
			"catering_instagram":  catering.Instagram,
			"catering_rating":     catering.Rating,
			"catering_priceRange": catering.PriceRange,
			"menus":               menus,
		}))
	})

	// get catering menu by id
	r.GET("/catering/:catering_id/cateringMenu/:menu_id", Middleware.Authorization(), func(c *gin.Context) {
		CateringID := c.Param("catering_id")
		id := c.Param("menu_id")

		var menu Model.CateringMenu
		if res := db.Where("catering_id = ?", CateringID).Where("id = ?", id).First(&menu); res.Error != nil {
			c.JSON(http.StatusNotFound, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry menu by id", gin.H{
			"menu": menu,
		}))
	})

	// get detailed catering menu by id
	r.GET("/catering/:catering_id/cateringMenu/:menu_id/detailed", Middleware.Authorization(), func(c *gin.Context) {
		CateringID := c.Param("catering_id")
		id := c.Param("menu_id")

		var detailedMenu Model.CateringMenuDetailed
		if res := db.Where("catering_id = ?", CateringID).Where("id = ?", id).First(&detailedMenu); res.Error != nil {
			c.JSON(http.StatusNotFound, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry menu by id", gin.H{
			"menu": detailedMenu,
		}))
	})
}
