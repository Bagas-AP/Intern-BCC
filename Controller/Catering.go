package Controller

// import (
// 	"bcc/Middleware"
// 	"bcc/Model"
// 	"bcc/Utils"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// func UserCatering(db *gorm.DB, q *gin.Engine) {
// 	r := q.Group("/api/user")

// 	// get all catering
// 	r.GET("allCatering", Middleware.Authorization(), func(c *gin.Context) {
// 		var caterings []Model.Catering
// 		if err := db.Find(&caterings).Error; err != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		Utils.HttpRespSuccess(c, http.StatusOK, "Queried all catering", caterings)
// 	})

// 	// search catering by name
// 	r.POST("/searchCatering", Middleware.Authorization(), func(c *gin.Context) {
// 		var input Model.CateringSearch
// 		if err := c.BindJSON(&input); err != nil {
// 			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
// 			return
// 		}

// 		var caterings []Model.Catering

// 		if res := db.Where("name LIKE ?", "%"+input.Name+"%").Find(&caterings); res.Error != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
// 			return
// 		}

// 		Utils.HttpRespSuccess(c, http.StatusOK, "Queried caterings by name", caterings)
// 	})

// 	// get catering by id and all its menu
// 	r.GET("/catering/:catering_id", Middleware.Authorization(), func(c *gin.Context) {
// 		id := c.Param("catering_id")

// 		var catering Model.Catering
// 		if res := db.Where("id = ?", id).First(&catering); res.Error != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
// 			return
// 		}

// 		var menus []Model.CateringMenu
// 		if res := db.Where("catering_id = ?", catering.ID).Find(&menus); res.Error != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
// 			return
// 		}

// 		Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering by id and it menus", gin.H{
// 			"catering_name":       catering.Name,
// 			"catering_address":    catering.Address,
// 			"catering_phone":      catering.Phone,
// 			"catering_instagram":  catering.Instagram,
// 			"catering_rating":     catering.Rating,
// 			"catering_priceRange": catering.PriceRange,
// 			"menus":               menus,
// 		})
// 	})

// 	// get catering menu by id
// 	r.GET("/catering/:catering_id/cateringMenu/:menu_index", Middleware.Authorization(), func(c *gin.Context) {
// 		CateringID := c.Param("catering_id")
// 		menuIndex, err := Utils.ParseStrToUint(c.Param("menu_index"))
// 		if err != nil {
// 			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		var menu Model.CateringMenu
// 		if res := db.Where("catering_id = ?", CateringID).Where("menu_index = ?", menuIndex).First(&menu); res.Error != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
// 			return
// 		}

// 		Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering menu by id", gin.H{
// 			"menu": menu,
// 		})
// 	})

// 	// get detailed catering menu by id
// 	r.GET("/catering/:catering_id/cateringMenu/:menu_index/detailed", Middleware.Authorization(), func(c *gin.Context) {
// 		CateringID := c.Param("catering_id")
// 		menuIndex, err := Utils.ParseStrToUint(c.Param("menu_index"))
// 		if err != nil {
// 			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		var detailedMenu Model.CateringMenuDetailed
// 		if res := db.Where("catering_id = ?", CateringID).Where("catering_menu_id = ?", menuIndex).First(&detailedMenu); res.Error != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
// 			return
// 		}

// 		Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering menu by id", gin.H{
// 			"menu": detailedMenu,
// 		})
// 	})

// 	// get catering tag hemat
// 	r.GET("/searchCateringTagHemat", Middleware.Authorization(), func(c *gin.Context) {
// 		var caterings []Model.CateringTags
// 		if res := db.Where("tag = ?", "Hemat").Preload("Catering.Seller").Find(&caterings); res.Error != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
// 			return
// 		}

// 		Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering menu by id", gin.H{
// 			"menu": caterings,
// 		})
// 	})

// 	// get catering tag vegetarian
// 	r.GET("/searchCateringTagVegetarian", Middleware.Authorization(), func(c *gin.Context) {
// 		var caterings []Model.CateringTags
// 		if res := db.Where("tag = ?", "Vegetarian").Preload("Catering.Seller").Find(&caterings); res.Error != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
// 			return
// 		}

// 		Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering menu by id", gin.H{
// 			"menu": caterings,
// 		})
// 	})

// 	// get catering tag diet
// 	r.GET("/searchCateringTagDiet", Middleware.Authorization(), func(c *gin.Context) {
// 		var caterings []Model.CateringTags
// 		if res := db.Where("tag = ?", "Diet").Preload("Catering.Seller").Find(&caterings); res.Error != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
// 			return
// 		}

// 		Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering menu by id", gin.H{
// 			"menu": caterings,
// 		})
// 	})

// 	// get catering tag sehat
// 	r.GET("/searchCateringTagSehat", Middleware.Authorization(), func(c *gin.Context) {
// 		var caterings []Model.CateringTags
// 		if res := db.Where("tag = ?", "Sehat").Preload("Catering.Seller").Find(&caterings); res.Error != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
// 			return
// 		}

// 		Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering menu by id", gin.H{
// 			"menu": caterings,
// 		})
// 	})

// 	// get catering tag halal
// 	r.GET("/searchCateringTagHalal", Middleware.Authorization(), func(c *gin.Context) {
// 		var caterings []Model.CateringTags
// 		if res := db.Where("tag = ?", "Halal").Preload("Catering.Seller").Find(&caterings); res.Error != nil {
// 			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
// 			return
// 		}

// 		Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering menu by id", gin.H{
// 			"menu": caterings,
// 		})
// 	})
// }
