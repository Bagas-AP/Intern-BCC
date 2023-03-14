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

}
