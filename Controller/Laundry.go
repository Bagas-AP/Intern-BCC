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

	// search laundry by name
	r.POST("/searchLaundry", Middleware.Authorization(), func(c *gin.Context) {
		var input Model.LaundrySearch

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusUnprocessableEntity, Utils.FailedResponse(err.Error()))
			return
		}

		var laundries []Model.Laundry

		if res := db.Where("name LIKE ?", "%"+input.Name+"%").Find(&laundries); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by name", gin.H{
			"data": laundries,
		}))
	})

	// get laundry by id and all its menu
	r.GET("/laundry/:laundry_id", Middleware.Authorization(), func(c *gin.Context) {
		id := c.Param("laundry_id")

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

	// get laundry menu by id
	r.GET("/laundry/:laundry_id/laundryMenu/:menu_index", Middleware.Authorization(), func(c *gin.Context) {
		laundryID := c.Param("laundry_id")
		menuIndex, err := Utils.ParseStrToUint(c.Param("menu_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(err.Error()))
			return
		}

		var menu Model.LaundryMenu
		if res := db.Where("laundry_id = ?", laundryID).Where("menu_index = ?", menuIndex).First(&menu); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry menu by id", gin.H{
			"data": menu,
		}))
	})

	// search laundry by tag Cepat
	r.GET("/searchLaundryTagCepat", Middleware.Authorization(), func(c *gin.Context) {
		var laundries []Model.LaundryTags
		if res := db.Where("tag = ?", "Cepat").Preload("Laundry.Seller").Find(&laundries); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by tags", gin.H{
			"data": laundries,
		}))
	})

	// search laundry by tag Baju
	r.GET("/searchLaundryTagBaju", Middleware.Authorization(), func(c *gin.Context) {
		var laundries []Model.LaundryTags
		if res := db.Where("tag = ?", "Baju").Preload("Laundry.Seller").Find(&laundries); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by tags", gin.H{
			"data": laundries,
		}))
	})

	// search laundry by tag Sepatu
	r.GET("/searchLaundryTagSepatu", Middleware.Authorization(), func(c *gin.Context) {
		var laundries []Model.LaundryTags
		if res := db.Where("tag = ?", "Sepatu").Preload("Laundry.Seller").Find(&laundries); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by tags", gin.H{
			"data": laundries,
		}))
	})

	// search laundry by tag Tas
	r.GET("/searchLaundryTagTas", Middleware.Authorization(), func(c *gin.Context) {
		var laundries []Model.LaundryTags
		if res := db.Where("tag = ?", "Tas").Preload("Laundry.Seller").Find(&laundries); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by tags", gin.H{
			"data": laundries,
		}))
	})

	// search laundry by tag Karpet
	r.GET("/searchLaundryTagKarpet", Middleware.Authorization(), func(c *gin.Context) {
		var laundries []Model.LaundryTags
		if res := db.Where("tag = ?", "Karpet").Preload("Laundry.Seller").Find(&laundries); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by tags", gin.H{
			"data": laundries,
		}))
	})

	// search laundry by tag Setrika
	r.GET("/searchLaundryTagSetrika", Middleware.Authorization(), func(c *gin.Context) {
		var laundries []Model.LaundryTags
		if res := db.Where("tag = ?", "Setrika").Preload("Laundry.Seller").Find(&laundries); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry by tags", gin.H{
			"data": laundries,
		}))
	})

	// add laundry order
	r.POST("/laundry/:laundry_id/laundryMenu/:menu_id", Middleware.Authorization(), func(c *gin.Context) {
		laundryID := c.Param("laundry_id")
		id, isIdExists := c.Params.Get("menu_id")

		if !isIdExists {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse("id not found"))
			return
		}

		var menu Model.LaundryMenu
		if res := db.Where("laundry_id = ?", laundryID).Where("id = ?", id).First(&menu); res.Error != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(res.Error.Error()))
			return
		}

		var input Model.InputLaundryOrder
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, Utils.FailedResponse(err.Error()))
			return
		}

		input = Model.InputLaundryOrder{
			LaundryMenuID: int(menu.ID),
			//LaundryID:     int(Utils.UintToString(id)),
			Quantity: input.Quantity,
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Success get laundry menu by id", gin.H{
			"data": input,
		}))
	})
}
