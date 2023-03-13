package Controller

import (
	"bcc/Config"
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func UserTransaction(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/transaction")

	r.POST("/catering/:catering_id/cateringMenu/:menu_index/transaction", Middleware.Authorization(), func(c *gin.Context) {
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

		var user Model.User
		ID, _ := c.Get("id")
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var menu Model.CateringMenu
		if res := db.Where("catering_id = ?", cateringID).Where("menu_index = ?", menuIndex).First(&menu); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var input Model.TransactionInputQuantity
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
		}

		var existingTransaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("menu_id = ?", menuIndex).Where("user_id = ?", user.ID).First(&existingTransaction); res.Error == nil {
			Utils.HttpRespFailed(c, http.StatusConflict, "transaction already exists")
			return
		}

		transaction := Model.Transaction{
			ID:        uuid.New(),
			UserID:    user.ID,
			Model:     2,
			ServiceID: int(cateringID),
			MenuID:    int(menuIndex),
			Quantity:  input.Quantity,
			Total:     (menu.Price * input.Quantity) + ((menu.Price * input.Quantity) * 10 / 100),
			CreatedAt: time.Now(),
		}

		if res := db.Create(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, res.Error.Error())
		}

		Utils.HttpRespSuccess(c, http.StatusCreated, "created new transaction", transaction)
	})

	r.GET("/catering/:catering_id/cateringMenu/:menu_index/transactionDetailed", Middleware.Authorization(), func(c *gin.Context) {
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

		var user Model.User
		ID, _ := c.Get("id")
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("model = ?", 2).Where("menu_id = ?", menuIndex).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
		}

		transactionResult := Model.TransactionResult{
			Name:            user.Name,
			Phone:           user.Phone,
			Address:         user.Address,
			MenuName:        menu.Name,
			MenuPrice:       menu.Price * transaction.Quantity,
			MenuExtendedFee: menu.Price * transaction.Quantity * 10 / 100,
			Total:           (menu.Price * transaction.Quantity) + ((menu.Price * transaction.Quantity) * 10 / 100),
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "transaction detail", transactionResult)
	})

	r.POST("/catering/:catering_id/cateringMenu/:menu_index/transactionDetailed", Middleware.Authorization(), func(c *gin.Context) {
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

		var user Model.User
		ID, _ := c.Get("id")
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("model = ?", 2).Where("menu_id = ?", menuIndex).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
		}

		var input Model.TransactionInputPayment
		if err := c.ShouldBindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		transaction.PaymentMethod = input.PaymentMethod

		if err := db.Save(&transaction).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "transaction detail", transaction)
	})

	r.GET("/catering/:catering_id/cateringMenu/:menu_index/payment", Middleware.Authorization(), func(c *gin.Context) {
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

		var user Model.User
		ID, _ := c.Get("id")
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("model = ?", 2).Where("menu_id = ?", menuIndex).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
		}

		transactionPayment := Model.TransactionPayment{
			PaymentDeadline: transaction.CreatedAt,
			Total:           transaction.Total,
		}

		if transaction.PaymentMethod == 1 {
			transactionPayment.AccountNumber = "BCA 2123982139"
			transactionPayment.AccountName = "Messi"
		} else if transaction.PaymentMethod == 2 {
			transactionPayment.AccountNumber = "2123982139"
			transactionPayment.AccountName = "Messi"
		} else if transaction.PaymentMethod == 3 {
			transactionPayment.AccountNumber = "2123982139"
			transactionPayment.AccountName = "Messi"
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "transaction detail", transactionPayment)
	})

	r.POST("/catering/:catering_id/cateringMenu/:menu_index/completePayment", Middleware.Authorization(), func(c *gin.Context) {
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

		var user Model.User
		ID, _ := c.Get("id")
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("model = ?", 2).Where("menu_id = ?", menuIndex).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
		}

		paymentProof, err := c.FormFile("payment_proof")
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		SupaBaseClient := Config.MakeSupaBaseClient()

		link, err := SupaBaseClient.Upload(paymentProof)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		transaction.PaymentProof = link
		transaction.Notes = c.PostForm("notes")

		if err := db.Save(&transaction).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "transaction detail", transaction)
	})

	r.POST("/catering/:catering_id/cateringMenu/:menu_index/uploadPaymentProof", Middleware.Authorization(), func(c *gin.Context) {
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

		var user Model.User
		ID, _ := c.Get("id")
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("model = ?", 2).Where("menu_id = ?", menuIndex).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
		}

		if err := db.Save(&transaction).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "transaction detail", transaction)
	})
}
