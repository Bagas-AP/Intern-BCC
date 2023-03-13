package Controller

import (
	"bcc/Middleware"
	"bcc/Model"
	"bcc/Utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

func UserTransaction(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/transaction")

	// add catering menu by id to transaction
	r.POST("/catering/:catering_id/cateringMenu/:menu_index/payment", Middleware.Authorization(), func(c *gin.Context) {
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

		var input Model.TransactionInput

		if err := c.BindJSON(&input); err != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
		}

		transactionResult := Model.TransactionResult{
			Name:      user.Name,
			Address:   user.Address,
			Phone:     user.Phone,
			MenuName:  menu.Name,
			MenuPrice: string(rune(menu.Price * input.Quantity)),
			Total:     (menu.Price * input.Quantity) + ((menu.Price * input.Quantity) * 10 / 100),
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "res", transactionResult)

		transaction := Model.Transaction{
			ID:            uuid.New(),
			UserID:        user.ID,
			Model:         2,
			ServiceID:     int(cateringID),
			MenuID:        int(menuIndex),
			Quantity:      input.Quantity,
			Total:         (menu.Price * input.Quantity) + ((menu.Price * input.Quantity) * 10 / 100),
			PaymentMethod: input.PaymentMethod,
		}

		var isExistTransaction Model.Transaction
		if err := db.Where("user_id = ?", user.ID).
			Where("model = ?", transaction.Model).
			Where("service_id = ?", transaction.ServiceID).
			Where("menu_id = ?", transaction.MenuID).
			First(&isExistTransaction).Error; err == nil {
			Utils.HttpRespFailed(c, http.StatusBadRequest, "Transaction already exist")
			return
		}

		if res := db.Create(&transaction); res.Error != nil {
			Utils.HttpRespFailed(c, http.StatusUnprocessableEntity, res.Error.Error())
		}

		Utils.HttpRespSuccess(c, http.StatusOK, "res", transaction)
	})
}
