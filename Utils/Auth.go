package Utils

import (
	"bcc/Model"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"gorm.io/gorm"
)

func SendToEmail(db *gorm.DB, email string, password string) error {
	var user Model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	auth := smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("EMAIL_PASS_APP"), "smtp.gmail.com")

	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Password Reset\r\n\r\n"+
			"Dear User,\r\n\r\n"+
			"We have received a request to reset your password for your account associated with this email. Please find below the details of the new password.\r\n\r\n"+
			"Phone: %s\r\n"+
			"Email: %s\r\n"+
			"New Password: %s\r\n\r\n"+
			"If you did not make this request, please contact our support team immediately.\r\n\r\n"+
			"Best regards,\r\n"+
			"%s",
		os.Getenv("EMAIL"), email, user.Phone, email, password, "TriServe",
	)

	err := smtp.SendMail("smtp.gmail.com"+":"+"587", auth, os.Getenv("EMAIL_FROM"), []string{email}, []byte(message))
	return err
}

func GenerateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 10
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}
