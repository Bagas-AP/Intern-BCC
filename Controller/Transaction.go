package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"net/http"
	"time"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserTransaction(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/transaction")

	// add new transaction
	r.POST("/:model_id/:service_id/:menu_id/new-transaction", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var user Model.User
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var input Model.TransactionInputQuantity
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		var menuPrice int

		if modelID == 1 {
			var menu Model.LaundryMenu
			if res := db.Where("laundry_id = ?", serviceID).Where("menu_index = ?", menuID).First(&menu); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			menuPrice = menu.Price
		} else if modelID == 2 {
			var menu Model.CateringMenu
			if res := db.Where("catering_id = ?", serviceID).Where("menu_index = ?", menuID).First(&menu); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			menuPrice = menu.Price
		}

		var existingTransaction Model.Transaction
		if res := db.Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).Where("user_id = ?", user.ID).First(&existingTransaction); res.Error == nil {
			Utils.HttpRespFailed(c, http.StatusConflict, "transaction already exists")
			return
		}

		transaction := Model.Transaction{
			ID:              uuid.New(),
			UserID:          user.ID,
			UserCity:        user.City,
			UserProvince:    user.Province,
			UserSubdistrict: user.Subdistrict,
			UserAddress:     user.Address,
			Model:           modelID,
			ServiceID:       int(serviceID),
			MenuID:          int(menuID),
			Quantity:        input.Quantity,
			Fee:             menuPrice * input.Quantity * 10 / 100,
			Total:           (menuPrice * input.Quantity) + ((menuPrice * input.Quantity) * 10 / 100),
			CreatedAt:       time.Now(),
		}

		if res := db.Create(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, res.Error.Error())
		}

		Utils.HttpRespSuccess(c, http.StatusCreated, "created new transaction", transaction)
	})

	// user change address
	r.POST("/:model_id/:service_id/:menu_id/change-address", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).Where("user_id = ?", ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var input Model.TransactionChangeAddress
		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		transaction.UserCity = input.City
		transaction.UserProvince = input.Province
		transaction.UserSubdistrict = input.Subdistrict
		transaction.UserAddress = input.Address

		if res := db.Save(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, res.Error.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "success change address", transaction)
	})

	// get transaction
	r.GET("/:model_id/:service_id/:menu_id/transaction-detailed", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var menuPrice int
		var menuName string

		if modelID == 1 {
			var menu Model.LaundryMenu
			if res := db.Where("laundry_id = ?", serviceID).Where("menu_index = ?", menuID).First(&menu); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			menuPrice = menu.Price
			menuName = menu.Name

		} else if modelID == 2 {
			var menu Model.CateringMenu
			if res := db.Where("catering_id = ?", serviceID).Where("menu_index = ?", menuID).First(&menu); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}

			menuPrice = menu.Price
			menuName = menu.Name
		}

		var user Model.User
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		transactionResult := Model.TransactionResult{
			Name:            user.Name,
			Phone:           user.Phone,
			City:            transaction.UserCity,
			Province:        transaction.UserProvince,
			Subdistrict:     transaction.UserSubdistrict,
			Address:         transaction.UserAddress,
			MenuName:        menuName,
			MenuPrice:       menuPrice * transaction.Quantity,
			MenuExtendedFee: menuPrice * transaction.Quantity * 10 / 100,
			Total:           (menuPrice * transaction.Quantity) + ((menuPrice * transaction.Quantity) * 10 / 100),
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "success get transaction detail", transactionResult)
	})

	// add payment method
	r.POST("/:model_id/:service_id/:menu_id/add-payment-method", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		if modelID == 1 {
			var menu Model.LaundryMenu
			if res := db.Where("laundry_id = ?", serviceID).Where("menu_index = ?", menuID).First(&menu); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}
		} else if modelID == 2 {
			var menu Model.CateringMenu
			if res := db.Where("catering_id = ?", serviceID).Where("menu_index = ?", menuID).First(&menu); res.Error != nil {
				Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
				return
			}
		}

		var user Model.User
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var input Model.TransactionInputPayment
		if err := c.BindJSON(&input); err != nil {
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

	// get payment method
	r.GET("/:model_id/:service_id/:menu_id/payment", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var user Model.User
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).Where("user_id = ?", user.ID).First(&transaction); res.Error != nil {
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

	// complete payment
	r.POST("/:model_id/:service_id/:menu_id//complete-payment", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).Where("user_id = ?", ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		paymentProof, err := c.FormFile("payment_proof")
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		SupaBaseClient := supabasestorageuploader.NewSupabaseClient(
			"https://flldkbhntqqaiflpxlhg.supabase.co",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImZsbGRrYmhudHFxYWlmbHB4bGhnIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY3NzU4Njk4OCwiZXhwIjoxOTkzMTYyOTg4fQ.CezKv4eOdEOyPEnVCqp3i0rNRLpz4MJOgL2GvM74QtQ",
			"photo",
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

	// cancel payment
	r.DELETE("/:model_id/:service_id/:menu_id/cancel-payment", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).Where("user_id = ?", ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		if err := db.Delete(&transaction).Error; err != nil {
			Utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "transaction deleted", nil)
	})

	// get detailed transaction
	r.GET("/:model_id/:service_id/:menu_id/confirm-payment", Middleware.Authorization(), func(c *gin.Context) {
		ID, modelID, serviceID, menuID, err := Utils.ParseIDs(c)
		if err != nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var transaction Model.Transaction
		if res := db.Where("model = ?", modelID).Where("service_id = ?", serviceID).Where("menu_id = ?", menuID).Where("user_id = ?", ID).First(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		newOrder := Model.Order{
			ID:        uuid.New(),
			UserID:    ID,
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
