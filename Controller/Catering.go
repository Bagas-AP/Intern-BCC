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

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get all catering", gin.H{
			"data": caterings,
		}))
	})

	// search catering by name
	r.POST("/searchCatering", Middleware.Authorization(), func(c *gin.Context) {
		var input Model.CateringSearch
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusUnprocessableEntity, Utils.FailedResponse(err.Error()))
			return
		}

		var caterings []Model.Catering

		if res := db.Where("name LIKE ?", "%"+input.Name+"%").Find(&caterings); res.Error != nil {
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
			"catering_name":       catering.Name,
			"catering_address":    catering.Address,
			"catering_phone":      catering.Phone,
			"catering_instagram":  catering.Instagram,
			"catering_rating":     catering.Rating,
			"catering_priceRange": catering.PriceRange,
			"menus":               menus,
		}))
	})

	// get catering menu by id
	r.GET("/catering/:catering_id/cateringMenu/:menu_index", Middleware.Authorization(), func(c *gin.Context) {
		CateringID := c.Param("catering_id")
		menuIndex, err := Utils.ParseStrToUint(c.Param("menu_index"))
		if err != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(err.Error()))
			return
		}

		var menu Model.CateringMenu
		if res := db.Where("catering_id = ?", CateringID).Where("id = ?", menuIndex).First(&menu); res.Error != nil {
			c.JSON(http.StatusNotFound, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry menu by id", gin.H{
			"menu": menu,
		}))
	})

	// get detailed catering menu by id
	r.GET("/catering/:catering_id/cateringMenu/:menu_index/detailed", Middleware.Authorization(), func(c *gin.Context) {
		CateringID := c.Param("catering_id")
		menuIndex, err := Utils.ParseStrToUint(c.Param("menu_index"))
		if err != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(err.Error()))
			return
		}

		var detailedMenu Model.CateringMenuDetailed
		if res := db.Where("catering_id = ?", CateringID).Where("catering_menu_id = ?", menuIndex).First(&detailedMenu); res.Error != nil {
			c.JSON(http.StatusNotFound, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry menu by id", gin.H{
			"menu": detailedMenu,
		}))
	})

	// hemat vegetarian diet sehat halal
	// get catering tag hemat
	r.GET("/searchLaundryTagHemat", Middleware.Authorization(), func(c *gin.Context) {
		var caterings []Model.CateringTags
		if res := db.Where("tag = ?", "Hemat").Preload("Catering.Seller").Find(&caterings); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get caterings by tags", gin.H{
			"data": caterings,
		}))
	})

	// get catering tag vegetarian
	r.GET("/searchLaundryTagVegetarian", Middleware.Authorization(), func(c *gin.Context) {
		var caterings []Model.CateringTags
		if res := db.Where("tag = ?", "Vegetarian").Preload("Catering.Seller").Find(&caterings); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get caterings by tags", gin.H{
			"data": caterings,
		}))
	})

	// get catering tag diet
	r.GET("/searchLaundryTagDiet", Middleware.Authorization(), func(c *gin.Context) {
		var caterings []Model.CateringTags
		if res := db.Where("tag = ?", "Diet").Preload("Catering.Seller").Find(&caterings); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get caterings by tags", gin.H{
			"data": caterings,
		}))
	})

	// get catering tag sehat
	r.GET("/searchLaundryTagSehat", Middleware.Authorization(), func(c *gin.Context) {
		var caterings []Model.CateringTags
		if res := db.Where("tag = ?", "Sehat").Preload("Catering.Seller").Find(&caterings); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get caterings by tags", gin.H{
			"data": caterings,
		}))
	})

	// get catering tag halal
	r.GET("/searchLaundryTagHalal", Middleware.Authorization(), func(c *gin.Context) {
		var caterings []Model.CateringTags
		if res := db.Where("tag = ?", "Halal").Preload("Catering.Seller").Find(&caterings); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get caterings by tags", gin.H{
			"data": caterings,
		}))
	})
}
