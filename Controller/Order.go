package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
					SellerPhone: laundry.Phone,
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
					SellerPhone: catering.Phone,
				})
			}
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Get all favourite", orderResults)
	})

	r.GET("/:model_id/:service_id/:menu_id/orderDetailed", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var order Model.Order
		if res := db.Where("user_id = ?", ID).Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).First(&order); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var laundry Model.Laundry
		var menuLaundry Model.LaundryMenu
		var catering Model.Catering
		var menuCatering Model.CateringMenu
		var orderRes Model.OrderResult
		var phone string

		if modelID == 1 {
			if res := db.Where("id = ?", order.ServiceID).First(&laundry); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			if res := db.Where("laundry_id = ?", order.ServiceID).Where("menu_index = ?", order.MenuID).First(&menuLaundry); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}
			phone = laundry.Phone
			orderRes.ModelID = modelID
			orderRes.ServiceID = order.ServiceID
			orderRes.MenuID = order.MenuID
			orderRes.ServiceName = laundry.Name
			orderRes.MenuName = menuLaundry.Name
			orderRes.MenuPhoto = menuLaundry.Photo
			orderRes.Status = order.Status
			orderRes.SellerPhone = phone

		} else if modelID == 2 {
			if res := db.Where("id = ?", order.ServiceID).First(&catering); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			if res := db.Where("catering_id = ?", order.ServiceID).Where("menu_index = ?", order.MenuID).First(&menuCatering); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}
			phone = catering.Phone
			orderRes.ModelID = order.Model
			orderRes.ServiceID = order.ServiceID
			orderRes.MenuID = order.MenuID
			orderRes.ServiceName = catering.Name
			orderRes.MenuName = menuCatering.Name
			orderRes.MenuPhoto = menuCatering.Photo
			orderRes.Status = order.Status
			orderRes.SellerPhone = phone
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "success get transaction detail", orderRes)
	})

	r.POST("/:model_id/:service_id/:menu_id/order-done", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var input Model.OrderInputConfirm
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if input.Confirm == 1 {
			var order Model.Order
			if res := db.Where("user_id = ?", ID).Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).First(&order); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

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

	r.POST("/:model_id/:service_id/:menu_id/rate", Middleware.Authorization(), func(c *gin.Context) {
		var input Model.InputRating
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		ID, modelID, serviceID, _, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var sumRating float64
		var count int64

		if err := db.Table("ratings").Where("model = ?", modelID).Where("service_id = ?", serviceID).
			Select("COALESCE(AVG(rating), 0) as rating, COUNT(*) as count").Row().Scan(&sumRating, &count); err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		rating := (sumRating + float64(input.Rating)) / float64(count+1)

		if modelID == 1 {
			var laundry Model.Laundry
			if res := db.Where("id = ?", serviceID).First(&laundry); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			}

			laundry.Rating = float64(rating)

			if res := db.Save(&laundry); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}
		} else if modelID == 2 {
			var catering Model.Catering
			if res := db.Where("id = ?", serviceID).First(&catering); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			}

			catering.Rating = float64(rating)

			if res := db.Save(&catering); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}
		}

		newRating := Model.Rating{
			UserID:    ID,
			Model:     modelID,
			ServiceID: int(serviceID),
			Rating:    float64(input.Rating),
		}

		if res := db.Create(&newRating); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "success rating", newRating)
	})
}
