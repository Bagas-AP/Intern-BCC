package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserFavourite(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/fav")

	// get all fav
	r.GET("/all", Middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		var favs []Model.Favourite
		if res := db.Where("user_id = ?", ID).Find(&favs); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var favResults []Model.FavouriteResult

		for _, fav := range favs {
			if fav.Model == 1 {
				var menu Model.LaundryMenu
				if res := db.Where("laundry_id = ?", fav.ServiceID).Where("menu_index = ?", fav.MenuID).First(&menu); res.Error != nil {
					Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
					return
				}
				favResults = append(favResults, Model.FavouriteResult{
					ModelID:   1,
					ServiceID: int(menu.LaundryID),
					MenuID:    menu.MenuIndex,
					MenuName:  menu.Name,
					MenuPrice: menu.Price,
					MenuPhoto: menu.Photo,
				})

			} else {
				var menu Model.CateringMenu
				if res := db.Where("catering_id = ?", fav.ServiceID).Where("menu_index = ?", fav.MenuID).First(&menu); res.Error != nil {
					Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
					return
				}
				favResults = append(favResults, Model.FavouriteResult{
					ModelID:   2,
					ServiceID: int(menu.CateringID),
					MenuID:    menu.MenuIndex,
					MenuName:  menu.Name,
					MenuPrice: menu.Price,
					MenuPhoto: menu.Photo,
				})
			}
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Get all favourite", favResults)
	})

	// add to fav menu by id
	r.POST("/:model_id/:service_id/:menu_id/add", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		favourite := Model.Favourite{
			UserID:    ID,
			Model:     modelID,
			ServiceID: int(serviceID),
			MenuID:    int(menuID),
		}

		if modelID == 1 {
			var isExistMenu Model.LaundryMenu
			if err := db.Where("laundry_id = ?", serviceID).Where("menu_index = ?", menuID).First(&isExistMenu).Error; err != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, "Menu doesn't exist")
				return
			}
		} else if modelID == 2 {
			var isExistMenu Model.CateringMenu
			if err := db.Where("catering_id = ?", serviceID).Where("menu_index = ?", menuID).First(&isExistMenu).Error; err != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, "Menu doesn't exist")
				return
			}
		}

		var isExistFavourite Model.Favourite
		if err := db.Where("user_id = ?", ID).Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).First(&isExistFavourite).Error; err == nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, "Menu already in favourite")
			return
		}

		if res := db.Create(&favourite); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Added new menu favourite", favourite)
	})

	// remove favourite menu by id
	r.DELETE("/:model_id/:service_id/:menu_id/delete", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		favourite := Model.Favourite{
			UserID:    ID,
			Model:     modelID,
			ServiceID: int(serviceID),
			MenuID:    int(menuID),
		}

		var isExistFavourite Model.Favourite
		if err := db.Where("user_id = ?", favourite.UserID).Where("model = ?", favourite.Model).Where("service_id = ?", favourite.ServiceID).Where("menu_id = ?", favourite.MenuID).First(&isExistFavourite).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, "favourite doesn't exist")
			return
		}

		if err := db.Where("user_id = ?", favourite.UserID).Where("model = ?", favourite.Model).Where("service_id = ?", favourite.ServiceID).Where("menu_id = ?", favourite.MenuID).Delete(&isExistFavourite).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "removed menu from favourite", nil)
	})
}
