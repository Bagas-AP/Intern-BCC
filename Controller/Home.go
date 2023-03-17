package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserHome(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/home")
	r.GET("/:model_id/all", Middleware.Authorization(), func(c *gin.Context) {
		modelID := c.Param("model_id")

		var laundries []Model.Laundry
		var caterings []Model.Catering
		if modelID == "1" {
			if err := db.Find(&laundries).Error; err != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusFound, "Queried all laundry", laundries)

		} else if modelID == "2" {
			if err := db.Find(&caterings).Error; err != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusFound, "Queried all catering", caterings)
		}
	})

	r.POST("/:model_id/search", Middleware.Authorization(), func(c *gin.Context) {
		var input Model.UserInput
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		modelID := c.Param("model_id")

		var laundries []Model.Laundry
		var caterings []Model.Catering
		if modelID == "1" {
			if res := db.Where("name LIKE ?", "%"+input.Name+"%").Find(&laundries); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusFound, "Queried all laundry", laundries)

		} else if modelID == "2" {
			if res := db.Where("name LIKE ?", "%"+input.Name+"%").Find(&caterings); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusFound, "Queried all catering", caterings)
		}
	})

	r.GET("/:model_id/:service_id", Middleware.Authorization(), func(c *gin.Context) {
		modelID := c.Param("model_id")
		serviceID := c.Param("service_id")

		if modelID == "1" {
			var laundry Model.Laundry
			if res := db.Where("id = ?", serviceID).First(&laundry); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			var menus []Model.LaundryMenu
			if res := db.Where("laundry_id = ?", laundry.ID).Find(&menus); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusOK, "Queried laundry by id", gin.H{
				"laundry_name":       laundry.Name,
				"laundry_address":    laundry.Address,
				"laundry_phone":      laundry.Phone,
				"laundry_instagram":  laundry.Instagram,
				"laundry_rating":     laundry.Rating,
				"laundry_priceRange": laundry.PriceRange,
				"menus":              menus,
			})

		} else if modelID == "2" {
			var catering Model.Catering
			if res := db.Where("id = ?", serviceID).First(&catering); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			var menus []Model.CateringMenu
			if res := db.Where("catering_id = ?", catering.ID).Find(&menus); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering by id and it menus", gin.H{
				"catering_name":       catering.Name,
				"catering_address":    catering.Address,
				"catering_phone":      catering.Phone,
				"catering_instagram":  catering.Instagram,
				"catering_rating":     catering.Rating,
				"catering_priceRange": catering.PriceRange,
				"menus":               menus,
			})
		}
	})

	r.GET("/:model_id/:service_id/:menu_id", Middleware.Authorization(), func(c *gin.Context) {
		modelID := c.Param("model_id")
		serviceID := c.Param("service_id")
		menuID := c.Param("menu_id")

		if modelID == "1" {
			var menu Model.LaundryMenu
			if res := db.Where("laundry_id = ?", serviceID).Where("menu_index = ?", menuID).First(&menu); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusOK, "Queried laundry menu by id", menu)

		} else if modelID == "2" {
			var menu Model.CateringMenu
			if res := db.Where("catering_id = ?", serviceID).Where("menu_index = ?", menuID).First(&menu); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering menu by id", menu)
		}
	})

	r.GET("/:model_id/:service_id/:menu_id/detailed", Middleware.Authorization(), func(c *gin.Context) {
		modelID := c.Param("model_id")
		serviceID := c.Param("service_id")
		menuID := c.Param("menu_id")

		if modelID == "1" {
			var menu Model.LaundryMenu
			if res := db.Where("laundry_id = ?", serviceID).Where("menu_index = ?", menuID).First(&menu); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusOK, "Queried laundry menu by id", menu)

		} else if modelID == "2" {
			var detailedMenu Model.CateringMenuDetailed
			if res := db.Where("catering_id = ?", serviceID).Where("catering_menu_id = ?", menuID).First(&detailedMenu); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusOK, "Queried detailed catering menu by id", detailedMenu)
		}
	})

	r.GET("/:model_id/tag/:tag", Middleware.Authorization(), func(c *gin.Context) {
		modelID := c.Param("model_id")
		tag := c.Param("tag")

		if modelID == "1" {
			var laundries []Model.LaundryTags
			if res := db.Where("tag = ?", tag).Preload("Laundry.Seller").Find(&laundries); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusOK, "Queried laundry by tags", laundries)
		} else if modelID == "2" {
			var caterings []Model.CateringTags
			if res := db.Where("tag = ?", tag).Preload("Catering.Seller").Find(&caterings); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusOK, "Queried catering by tags", caterings)
		}
	})
}
