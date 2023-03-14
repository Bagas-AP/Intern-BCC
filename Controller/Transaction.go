package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
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
			return
		}

		var existingTransaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("menu_id = ?", menuIndex).Where("user_id = ?", user.ID).First(&existingTransaction); res.Error == nil {
			Utils.HttpRespFailed(c, http.StatusConflict, "transaction already exists")
			return
		}

		transaction := Model.Transaction{
			ID:          uuid.New(),
			UserID:      user.ID,
			UserAddress: user.Province + user.City + user.Subdistrict + user.Address,
			Model:       2,
			ServiceID:   int(cateringID),
			MenuID:      int(menuIndex),
			Quantity:    input.Quantity,
			Total:       (menu.Price * input.Quantity) + ((menu.Price * input.Quantity) * 10 / 100),
			CreatedAt:   time.Now(),
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
			return
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

		Utils.HttpRespSuccess(c, http.StatusOK, "success get transaction detail", transactionResult)
	})

	r.POST("/catering/:catering_id/cateringMenu/:menu_index/changeAddress", Middleware.Authorization(), func(c *gin.Context) {
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

		var transaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("model = ?", 2).Where("menu_id = ?", menuIndex).Where("user_id = ?", ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var input Model.TransactionChangeAddress
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		transaction.UserAddress = input.Province + input.City + input.Subdistrict + input.Address

		if res := db.Save(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "success change address", transaction)
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
			return
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

		Utils.HttpRespSuccess(c, http.StatusOK, "successfully added transaction payment", transaction)
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
			return
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

		Utils.HttpRespSuccess(c, http.StatusOK, "detailed transaction", transactionPayment)
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

		var user Model.User
		ID, _ := c.Get("id")
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("model = ?", 2).Where("menu_id = ?", menuIndex).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		paymentProof, err := c.FormFile("payment_proof")
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		SupaBaseClient := supabasestorageuploader.NewSupabaseClient(
			"https://llghldcosbvakddztrpt.supabase.co",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImxsZ2hsZGNvc2J2YWtkZHp0cnB0Iiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY3NzUwMjEzOCwiZXhwIjoxOTkzMDc4MTM4fQ.1fHxjZfPYIaTEEzEh9aUg5_T0yh-cyq5KG9KAbxt8C4",
			"picture",
			"",
		)

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

		Utils.HttpRespSuccess(c, http.StatusOK, "added notes and payment proof", transaction)
	})

	r.DELETE("/catering/:catering_id/cateringMenu/:menu_index/cancelPayment", Middleware.Authorization(), func(c *gin.Context) {
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

		var transaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("model = ?", 2).Where("menu_id = ?", menuIndex).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		if err := db.Delete(&transaction).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "transaction deleted", nil)
	})

	r.POST("/catering/:catering_id/cateringMenu/:menu_index/confirmPayment", Middleware.Authorization(), func(c *gin.Context) {
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

		var transaction Model.Transaction
		if res := db.Where("service_id = ?", cateringID).Where("model = ?", 2).Where("menu_id = ?", menuIndex).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		newOrder := Model.Order{
			ID:        uuid.New(),
			UserID:    user.ID,
			Model:     transaction.Model,
			ServiceID: transaction.ServiceID,
			MenuID:    transaction.MenuID,
			Status:    1,
			CreatedAt: time.Now(),
		}

		if err := db.Create(&newOrder).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		if err := db.Delete(&transaction).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "order created", newOrder)
	})
}
