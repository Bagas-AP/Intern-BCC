package Controller

import (
	"bcc/Model"
	"bcc/Utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func UserFinance(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/finance")
	r.GET("/all", func(c *gin.Context) {
		var categories []Model.WalletCategories
		if res := db.Find(&categories); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "queried wallet categories", categories)
	})

	r.GET("/category/:id", func(c *gin.Context) {
		categoryID := c.Param("id")

		var category Model.WalletCategories
		if res := db.Where("id = ?", categoryID).First(&category); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "queried wallet category", category)
	})
}
