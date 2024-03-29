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

func UserProfile(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")

	// see user profile
	r.GET("/profile", Middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var user Model.User
		if err := db.Where("id = ?", ID).Take(&user).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Profile fetched successfully", gin.H{
			"photo": user.ProfilePicture,
			"nama":  user.Name,
			"phone": user.Phone,
			"email": user.Email,
			// "user": user,
		})
	})

	// edit user profile
	r.PATCH("/profile", Middleware.Authorization(), func(c *gin.Context) {
		var input Model.UserUpdate
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		ID, _ := c.Get("id")

		user := Model.User{
			Name:           input.Name,
			Phone:          input.Phone,
			Email:          input.Email,
			Password:       Utils.Hash(input.Password),
			Province:       input.Province,
			City:           input.City,
			Subdistrict:    input.Subdistrict,
			Address:        input.Address,
			ProfilePicture: input.ProfilePicture,
			UpdatedAt:      time.Now(),
		}

		if err := db.Where("id = ?", ID).Model(&user).Updates(user).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Profile updated successfully", user)
	})

	// user create a support tiket
	r.POST("/user/support", Middleware.Authorization(), func(c *gin.Context) {
		ID, _, _, _, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var input Model.Support
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
		}

		support := Model.Support{
			UserID:    ID,
			Help:      input.Help,
			CreatedAt: time.Now(),
		}

		if err := db.Create(&support).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "Support created successfully", support)
	})
}
