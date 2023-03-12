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

func Favourite(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/fav")

	// get all fav
	r.GET("/all", Middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var favourites []Model.Favourite
		if res := db.Where("user_id = ?", ID.(uint)).Find(&favourites); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Queried all favourite", favourites)
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

		if res := db.Create(&favourite); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Added new catering favourite", favourite)
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

		if res := db.Create(&favourite); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Added new laundry favourite", favourite)
	})
}
