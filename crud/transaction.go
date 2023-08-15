package crud

import (
	"fmt"

	model "github.com/Prokuma/ProkumaLabAccount-Backend/models"
)

func CreateBook(user *model.User, book *model.Book) error {
	tx := DB.Begin()
	result := tx.Create(book)

	if result.Error != nil {
		fmt.Println("Book could not create: ", result.Error)
		tx.Rollback()
		return result.Error
	}

	err := tx.Create(&model.BookAuthorization{
		UserId:    *&user.UserId,
		BookId:    *&book.BookId,
		Authority: "admin,read,write,update,delete",
	}).Error

	if err != nil {
		fmt.Println("Authorization could not create: ", err)
		tx.Rollback()
		return err
	}

	err = tx.Create(&model.AccountTitle{
		BookId: *&book.BookId,
		Name:   "omission_system",
		Type:   4,
	}).Error

	if err != nil {
		fmt.Println("Authorization could not create: ", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error

	if err != nil {
		return err
	}

	return nil
}

func UpdateBook(book *model.Book) error {
	err := DB.Save(book).Error

	if err != nil {
		fmt.Println("Book could not update: ", err)
		return err
	}

	return nil
}

func DeleteBook(bookId string) error {
	err := DB.Unscoped().Delete(&model.Book{BookId: bookId}).Error

	if err != nil {
		fmt.Println("Delete the book was failed: " + err.Error())
		return err
	}

	return nil
}

func GetBook(bookId string) (model.Book, error) {
	var book model.Book
	err := DB.Where(&model.Book{BookId: bookId}).First(&book).Error

	if err != nil {
		fmt.Println("Book could not found: ", err)
		return model.Book{}, err
	}

	return book, nil
}

func CreateAccountTitle(title *model.AccountTitle) error {
	err := DB.Create(title).Error

	if err != nil {
		fmt.Println("Account Title could not create: ", err)
		return err
	}

	return nil
}

func GetAccountTitle(book *model.Book, accountTitleId uint64) (model.AccountTitle, error) {
	var accountTitle model.AccountTitle
	err := DB.Where("book_id = ?", book.BookId).Find(&accountTitle).Error

	if err != nil {
		return model.AccountTitle{}, err
	}

	return accountTitle, nil
}

func DeleteAccountTitle(book *model.Book, accountTitleId uint64) error {
	err := DB.Delete(&model.AccountTitle{AccountTitleId: accountTitleId}).Error

	if err != nil {
		return err
	}

	return nil
}

func UpdateAccountTitle(title *model.AccountTitle) error {
	err := DB.Save(title).Error

	if err != nil {
		fmt.Println("Account Title could not create: ", err)
		return err
	}

	return nil
}

func CreateBookAuthorization(authorization *model.BookAuthorization) error {
	err := DB.Create(authorization).Error

	if err != nil {
		fmt.Println("Add authorized user was failed: ", err)
		return err
	}

	return nil
}

func GetBookAuthorization(user *model.User, book *model.Book) (model.BookAuthorization, error) {
	var bookAuthorization model.BookAuthorization
	err := DB.Where(&model.BookAuthorization{UserId: *&user.UserId}).First(&bookAuthorization).Error

	if err != nil {
		fmt.Println("Book Authorization not found: ", err)
		return model.BookAuthorization{}, err
	}

	return bookAuthorization, nil
}

func UpdateBookAuthorization(authorization *model.BookAuthorization) error {
	err := DB.Save(authorization).Error

	if err != nil {
		fmt.Println("Add authorized user was failed: ", err)
		return err
	}

	return nil
}

func DeleteBookAuthorization(authorization *model.BookAuthorization) error {
	err := DB.Delete(authorization).Error

	if err != nil {
		fmt.Println("Delete authorized user was failed: ", err)
		return err
	}

	return nil
}

func CreateBookAndAccountTitleFromBook(year uint, name string, admin *model.User, oldBook *model.Book) error {
	newBook := model.Book{Year: year, Name: name}
	tx := DB.Begin()
	result := tx.Create(&newBook)

	if result.Error != nil {
		tx.Rollback()
		fmt.Println("Create New Book was failed: ", result.Error)
		return result.Error
	}

	err := DB.Create(&model.BookAuthorization{
		BookId:    newBook.BookId,
		UserId:    *&admin.UserId,
		Authority: "admin,read,write,update,delete",
	}).Error

	var oldAccountTtiles []model.AccountTitle
	var newAccountTitles []model.AccountTitle
	tx.Where(&model.AccountTitle{BookId: *&oldBook.BookId}).Find(&oldAccountTtiles)
	for _, accountTitle := range oldAccountTtiles {
		newAccountTitles = append(newAccountTitles, model.AccountTitle{
			BookId: newBook.BookId,
			Name:   accountTitle.Name,
			Amount: accountTitle.Amount, // 繰越
			Type:   accountTitle.Type,
		})
	}
	err = tx.Create(&newAccountTitles).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Create New Account titles was failed: ", err)
		return err
	}

	err = tx.Commit().Error

	if err != nil {
		fmt.Println("Create Book from old Book was failed: ", err)
		return err
	}

	return nil
}

func CreateTransaction(transaction *model.Transaction, subTransactions *[]model.SubTransaction) error {
	tx := DB.Begin()
	err := tx.Create(transaction).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Transaction Create Error: ", err)
		return err
	}
	err = tx.Create(subTransactions).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Sub Transaction Create Error: ", err)
		return err
	}

	for _, subTransaction := range *subTransactions {
		var debit model.AccountTitle
		tx.Where(&model.AccountTitle{AccountTitleId: subTransaction.DebitId}).First(&debit)
		if debit.Type%2 == 0 {
			debit.Amount += subTransaction.Amount
		} else {
			debit.Amount -= subTransaction.Amount
		}

		var credit model.AccountTitle
		tx.Where(&model.AccountTitle{AccountTitleId: subTransaction.CreditId}).First(&credit)
		if credit.Type%2 == 0 {
			credit.Amount -= subTransaction.Amount
		} else {
			credit.Amount += subTransaction.Amount
		}

		err = tx.Save(debit).Error
		if err != nil {
			fmt.Println("Account Title Update Error: ", err)
			tx.Rollback()
			return err
		}
		err = tx.Save(credit).Error
		if err != nil {
			tx.Rollback()
			fmt.Println("Account Title Update Error: ", err)
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		fmt.Println("Transaction Commit Error: ", err)
		return err
	}

	return nil
}

func DeleteTransaction(transaction *model.Transaction) error {
	var subTransactions []model.SubTransaction

	tx := DB.Begin()
	tx.Where(&model.Transaction{TransactionId: transaction.TransactionId}).Find(&subTransactions)
	err := tx.Delete(transaction).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Transaction Delete Error: ", err)
		return err
	}

	for _, subTransaction := range subTransactions {
		var debit model.AccountTitle
		tx.Where(&model.AccountTitle{AccountTitleId: subTransaction.DebitId}).First(&debit)
		if debit.Type%2 == 0 {
			debit.Amount -= subTransaction.Amount
		} else {
			debit.Amount += subTransaction.Amount
		}

		var credit model.AccountTitle
		tx.Where(&model.AccountTitle{AccountTitleId: subTransaction.CreditId}).First(&credit)
		if credit.Type%2 == 0 {
			credit.Amount += subTransaction.Amount
		} else {
			credit.Amount -= subTransaction.Amount
		}

		err = tx.Save(debit).Error
		if err != nil {
			fmt.Println("Account Title Update Error: ", err)
			tx.Rollback()
			return err
		}
		err = tx.Save(credit).Error
		if err != nil {
			tx.Rollback()
			fmt.Println("Account Title Update Error: ", err)
			return err
		}
	}
	err = tx.Delete(&subTransactions).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Sub Transactions Delete Error: ", err)
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		fmt.Println("Transaction Commit Error: ", err)
		return err
	}

	return nil
}

func GetTransaction(book *model.Book, transactionId uint64) (model.Transaction, error) {
	var transaction model.Transaction
	err := DB.Where(&model.Transaction{TransactionId: transactionId, BookId: *&book.BookId}).First(&transaction).Error

	if err != nil {
		fmt.Println("No Transaction")
		return model.Transaction{}, err
	}

	return transaction, nil
}

func GetSubTransaction(book *model.Book, subTransactionId uint64) (model.SubTransaction, error) {
	var subTransaction model.SubTransaction
	err := DB.Where(&model.SubTransaction{SubTransactionId: subTransactionId, BookId: *&book.BookId}).First(&subTransaction).Error

	if err != nil {
		fmt.Println("No Sub Transaction")
		return model.SubTransaction{}, err
	}

	return subTransaction, nil
}

func GetTransactions(book *model.Book, dataPerPage int, page int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := DB.Where(&model.Transaction{BookId: *&book.BookId}).Offset(dataPerPage * page).Limit(page).Find(&transactions).Error

	if err != nil {
		fmt.Println("No Transactions")
		return nil, err
	}

	return transactions, nil
}

func GetSubTransactionsFromTransactionIds(book *model.Book, start uint64, end uint64) ([]model.SubTransaction, error) {
	var subTransactions []model.SubTransaction
	err := DB.Where("book_id = ? AND (transaction_id BETWEEN ? AND ?)", book.BookId, start, end).Find(&subTransactions).Error

	if err != nil {
		fmt.Println("No Sub Transactions")
		return nil, err
	}

	return subTransactions, nil
}
