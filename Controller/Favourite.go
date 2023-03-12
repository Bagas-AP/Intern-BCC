package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"net/http"
	"strconv"

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

			} else if fav.Model == 2 {
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

	// add to fav catering menu by id
	r.POST("/catering/:catering_id/cateringMenu/:menu_index", Middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
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

		favourite := Model.Favourite{
			UserID:    ID.(uint),
			Model:     2,
			ServiceID: int(cateringID),
			MenuID:    int(menuIndex),
		}

		var isExistFavourite Model.Favourite
		if err := db.Where("user_id = ?", ID).Where("model = ?", favourite.Model).Where("service_id = ?", favourite.ServiceID).Where("menu_id = ?", favourite.MenuID).First(&isExistFavourite).Error; err == nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, "Menu already in favourite")
			return
		}

		if res := db.Create(&favourite); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Added new catering menu favourite", favourite)
	})

	// add to fav laundry menu by id
	r.POST("/laundry/:laundry_id/laundryMenu/:menu_index", Middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		laundryID, err := strconv.Atoi(c.Param("laundry_id"))
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}
		menuIndex, err := strconv.Atoi(c.Param("menu_index"))
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		favourite := Model.Favourite{
			UserID:    ID.(uint),
			Model:     1,
			ServiceID: int(laundryID),
			MenuID:    int(menuIndex),
		}

		var isExistFavourite Model.Favourite
		if err := db.Where("user_id = ?", ID).Where("model = ?", favourite.Model).Where("service_id = ?", favourite.ServiceID).Where("menu_id = ?", favourite.MenuID).First(&isExistFavourite).Error; err == nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, "Menu already in favourite")
			return
		}

		if res := db.Create(&favourite); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Added new laundry menu favourite", favourite)
	})

	// remove favourite laundry menu by id
	r.DELETE("/laundry/:laundry_id/laundryMenu/:menu_index", Middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		laundryID, err := strconv.Atoi(c.Param("laundry_id"))
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}
		menuIndex, err := strconv.Atoi(c.Param("menu_index"))
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		favourite := Model.Favourite{
			UserID:    ID.(uint),
			Model:     1,
			ServiceID: laundryID,
			MenuID:    menuIndex,
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

		Utils.HttpRespSuccess(c, http.StatusOK, "removed laundry menu from favourite", "")
	})

	// remove favourite catering menu by id
	r.DELETE("/catering/:catering_id/cateringMenu/:menu_index", Middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
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

		favourite := Model.Favourite{
			UserID:    ID.(uint),
			Model:     2,
			ServiceID: int(cateringID),
			MenuID:    int(menuIndex),
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

		Utils.HttpRespSuccess(c, http.StatusOK, "removed catering menu from favourite", "")
	})
}
