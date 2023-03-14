package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func UserOrder(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/order")

	// get all order
	r.GET("/all", Middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		var orders []Model.Order
		if res := db.Where("user_id = ?", ID).Find(&orders); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var orderResults []Model.OrderResult

		for _, order := range orders {
			if order.Model == 1 {
				var laundry Model.Laundry
				if res := db.Where("id = ?", order.ServiceID).First(&laundry); res.Error != nil {
					Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
					return
				}

				var menu Model.LaundryMenu
				if res := db.Where("laundry_id = ?", order.ServiceID).Where("menu_index = ?", order.MenuID).First(&menu); res.Error != nil {
					Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
					return
				}
				orderResults = append(orderResults, Model.OrderResult{
					ModelID:     1,
					ServiceID:   int(menu.LaundryID),
					MenuID:      menu.MenuIndex,
					ServiceName: laundry.Name,
					MenuName:    menu.Name,
					MenuPhoto:   menu.Photo,
				})

			} else if order.Model == 2 {
				var catering Model.Catering
				if res := db.Where("id = ?", order.ServiceID).First(&catering); res.Error != nil {
					Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
					return
				}

				var menu Model.CateringMenu
				if res := db.Where("catering_id = ?", order.ServiceID).Where("menu_index = ?", order.MenuID).First(&menu); res.Error != nil {
					Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
					return
				}
				orderResults = append(orderResults, Model.OrderResult{
					ModelID:     2,
					ServiceID:   int(menu.CateringID),
					MenuID:      menu.MenuIndex,
					ServiceName: catering.Name,
					MenuName:    menu.Name,
					MenuPhoto:   menu.Photo,
				})
			}
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Get all favourite", orderResults)
	})

	r.GET("/catering/:catering_id/cateringMenu/:menu_index/orderDetailed", Middleware.Authorization(), func(c *gin.Context) {
		cateringID, err := Utils.ParseStrToUint(c.Param("catering_id"))
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}
		menuIndex, err := Utils.ParseStrToUint(c.Param("menu_index"))
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var menu Model.CateringMenu
		if res := db.Where("catering_id = ?", cateringID).Where("menu_index = ?", menuIndex).First(&menu); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		ID, _ := c.Get("id")

		var order Model.OrderResult
		if res := db.Where("user_id = ?", ID).Where("model = ?", 2).Where("service_id = ?", cateringID).Where("menu_id = ?", menuIndex).First(&order); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var laundry Model.Laundry
		var catering Model.Catering
		var phone string

		if order.ModelID == 1 {
			if res := db.Where("id = ?", order.ServiceID).First(&laundry); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}
			phone = laundry.Phone
		} else if order.ModelID == 2 {
			if res := db.Where("id = ?", order.ServiceID).First(&catering); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}
			phone = catering.Phone
		}

		orderResult := Model.OrderResult{
			ModelID:     order.ModelID,
			ServiceID:   order.ServiceID,
			MenuID:      order.MenuID,
			ServiceName: order.ServiceName,
			MenuName:    order.MenuName,
			MenuPhoto:   order.MenuPhoto,
			Status:      order.Status,
			SellerPhone: phone,
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "success get transaction detail", orderResult)
	})

	r.POST("/catering/:catering_id/cateringMenu/:menu_index/orderDone", Middleware.Authorization(), func(c *gin.Context) {
		cateringID, err := Utils.ParseStrToUint(c.Param("catering_id"))
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}
		menuIndex, err := Utils.ParseStrToUint(c.Param("menu_index"))
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		ID, _ := c.Get("id")

		var input Model.OrderInputConfirm
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if input.Confirm == 1 {
			var order Model.Order
			if res := db.Where("user_id = ?", ID).Where("model = ?", 2).Where("service_id = ?", cateringID).Where("menu_id = ?", menuIndex).First(&order); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}
			//order.CompletedAt = time.Now()

			//if res := db.Save(&order); res.Error != nil {
			//	Utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			//	return
			//}
			if res := db.Delete(&order); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}

			Utils.HttpRespSuccess(c, http.StatusOK, "order has been completed", nil)

		} else {
			Utils.HttpRespFailed(c, http.StatusBadRequest, "confirm must be 1")
			return
		}
	})

	r.POST("/catering/:catering_id/cateringMenu/:menu_index/rating", Middleware.Authorization(), func(c *gin.Context) {
		var input Model.OrderInputRating
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		cateringID, err := Utils.ParseStrToUint(c.Param("catering_id"))
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var sumRating int
		var count int

		if err := db.Table("ratings").Select("AVG(rating) as rating, COUNT(*) as count").Where("model = ? AND service_id = ?", 1, cateringID).Row().Scan(&sumRating, &count); err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		rating := (sumRating + input.Rating) / (count + 1)

		var catering Model.Catering
		if res := db.Where("id = ?", cateringID).First(&catering); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
		}

		catering.Rating = rating

		if res := db.Save(&catering); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "success rating", catering)
	})
}
