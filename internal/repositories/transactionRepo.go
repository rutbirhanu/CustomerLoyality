package repositories

type TransactionRepo interface {
	CreateEntry()
	UpdateBalance()
	CreateTransaction()
}

// func Transaction(amount float64, types string, action string) {

// 	switch action {
// 	case "debit":
// 		amount = -amount
// 	case "credit":
// 		amount = +amount
// 	}

// 	switch
// }
