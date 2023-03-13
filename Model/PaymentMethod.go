package Model

type AdminPaymentAccount struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	Bank       string `json:"bank"`
	BankNumber string `json:"bank_number"`
	OVO        string `json:"ovo"`
	Dana       string `json:"dana"`
}
