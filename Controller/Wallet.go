package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserWallet(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/finance")
	// get all wallet categories
	r.GET("/all", func(c *gin.Context) {
		var categories []Model.WalletCategories
		if res := db.Find(&categories); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "queried wallet categories", categories)
	})

	// create new wallet category
	r.POST("/new-category", Middleware.Authorization(), func(c *gin.Context) {
		var input Model.WalletCategoryInput
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		ID, _ := c.Get("id")
		var wallet Model.Wallet
		if res := db.Where("user_id = ?", ID).First(&wallet); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var maxIndex int
		err := db.Model(&Model.WalletCategories{}).Select("MAX(`index`)").Where("wallet_id = ?", wallet.ID).Scan(&maxIndex).Error
		if err != nil {
			log.Printf("error max index: %v", maxIndex)
			maxIndex = 0
		}

		newCategory := Model.WalletCategories{
			ID:       uuid.New(),
			UserID:   ID.(uint),
			WalletID: wallet.ID,
			Index:    maxIndex + 1,
			Name:     input.Name,
			Budget:   input.Budget,
			Balance:  input.Budget,
		}

		if res := db.Create(&newCategory); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "created new wallet category", newCategory)
	})

	// get wallet category by id and all its details
	r.GET("/category/:id", Middleware.Authorization(), func(c *gin.Context) {
		index := c.Param("id")
		ID, _ := c.Get("id")

		var category Model.WalletCategories
		if res := db.Preload("WalletTransactions").Where("index = ?", index).Where("user_id = ?", ID).First(&category); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "queried wallet category", category)
	})

	// add new transaction to wallet category by id
	r.POST("/category/:id/add-transaction", Middleware.Authorization(), func(c *gin.Context) {
		index := c.Param("id")
		ID, _ := c.Get("id")

		var category Model.WalletCategories
		if res := db.Preload("WalletTransactions").Where("index = ?", index).Where("user_id = ?", ID).First(&category); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var input Model.WalletTransactionInput
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		transaction := Model.WalletTransaction{
			ID:               uuid.New(),
			WalletID:         category.WalletID,
			UserID:           category.UserID,
			WalletCategoryID: category.ID,
			Expense:          input.Expense,
			Amount:           input.Amount,
			CreatedAt:        time.Time{},
		}

		if res := db.Create(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		category.Budget = category.Budget - input.Amount

		if res := db.Save(&category); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "added new expense transaction", transaction)
	})

	// delete wallet category by id
	r.DELETE("/category/:id/delete-category", Middleware.Authorization(), func(c *gin.Context) {
		index := c.Param("id")
		ID, _ := c.Get("id")

		var category Model.WalletCategories
		if res := db.Where("index = ?", index).Where("user_id = ?", ID).First(&category); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var transactions []Model.WalletTransaction
		if res := db.Where("wallet_id = ?", category.WalletID).Find(&transactions); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusForbidden, res.Error.Error())
			return
		}

		if len(transactions) > 0 {
			if err := db.Delete(&transactions); err.Error != nil {
				Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error.Error())
				return
			}
		}

		if err := db.Delete(&category); err.Error != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "deleted wallet category", category)
	})

}
