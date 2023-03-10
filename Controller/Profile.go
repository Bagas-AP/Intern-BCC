package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Profile(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")

	// see user profile
	r.GET("/profile", Middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var user Model.User
		if err := db.Where("id = ?", ID).Take(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, Utils.FailedResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Profile fetched successfully", gin.H{
			"success": true,
			// "data":    user,
			"photo": user.ProfilePicture,
			"nama":  user.Name,
			"phone": user.Phone,
			"email": user.Email,
		}))
	})

	// edit user profile
	r.PATCH("/profile", Middleware.Authorization(), func(c *gin.Context) {
		var input Model.UserUpdate
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusUnprocessableEntity, Utils.FailedResponse(err.Error()))
			return
		}

		ID, _ := c.Get("id")

		//var user Model.User
		//if err := db.Where("id = ?", ID).Take(&user); err != nil {
		//	c.JSON(http.StatusNotFound, Utils.FailedResponse(err.Error.Error()))
		//	return
		//}

		user := Model.User{
			Name:           input.Name,
			Phone:          input.Phone,
			Email:          input.Email,
			Password:       Hash(input.Password),
			Province:       input.Province,
			City:           input.City,
			Subdistrict:    input.Subdistrict,
			Address:        input.Address,
			ProfilePicture: input.ProfilePicture,
			UpdatedAt:      time.Now(),
		}

		if err := db.Where("id = ?", ID).Model(&user).Updates(user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Utils.FailedResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Utils.SucceededReponse("Profile updated successfully", gin.H{
			"success": true,
			"data":    user,
		}))

	})
}
