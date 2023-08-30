package crud

import (
	"fmt"

	model "github.com/Prokuma/PLAccounting-Backend/models"
	"gorm.io/gorm"
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

	err = tx.Commit().Error

	if err != nil {
		return err
	}

	return nil
}

func UpdateBook(book *model.Book) error {
	err := DB.Model(&model.Book{BookId: book.BookId}).Updates(book).Error

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

func GetBookPages(bookId string, dataPerPages int) int64 {
	var count int64
	err := DB.Model(&model.Transaction{}).Where(&model.Transaction{BookId: bookId}).Count(&count).Error

	if err != nil {
		fmt.Println("Book could not found: ", err)
		return 1
	}

	if (count % int64(dataPerPages)) == 1 {
		return count/int64(dataPerPages) + 1
	}
	return count / int64(dataPerPages)
}

func GetAllBooks(user *model.User) (*[]model.Book, error) {
	var books []model.Book
	subQuery := DB.Select("book_id").Where(&model.BookAuthorization{UserId: user.UserId}).Table("book_authorizations")
	err := DB.Having("book_id = (?)", subQuery).Find(&books).Error

	if err != nil {
		fmt.Println("Books could not found: ", err)
		return nil, err
	}

	return &books, nil
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
	err := DB.Where(&model.AccountTitle{AccountTitleId: accountTitleId, BookId: book.BookId}).First(&accountTitle).Error

	if err != nil {
		return model.AccountTitle{}, err
	}

	return accountTitle, nil
}

func GetAllAccountTitles(book *model.Book) (*[]model.AccountTitle, error) {
	var accountTitles []model.AccountTitle

	err := DB.Where(&model.AccountTitle{BookId: *&book.BookId}).Order("account_title_id DESC").Find(&accountTitles).Error

	if err != nil {
		fmt.Println("No Account Titles", err)
		return nil, err
	}

	return &accountTitles, nil
}

func GetAccountTitles(book *model.Book, dataPerPage int, page int) (*[]model.AccountTitle, error) {
	var accountTitles []model.AccountTitle

	err := DB.Where(&model.AccountTitle{BookId: *&book.BookId}).Order("account_title_id DESC").Limit(dataPerPage).Offset(dataPerPage * page).Find(&accountTitles).Error

	if err != nil {
		fmt.Println("No Account Titles", err)
		return nil, err
	}

	return &accountTitles, nil
}

func DeleteAccountTitle(book *model.Book, accountTitleId uint64) error {
	err := DB.Delete(&model.AccountTitle{AccountTitleId: accountTitleId}).Error

	if err != nil {
		return err
	}

	return nil
}

func UpdateAccountTitle(title *model.AccountTitle) error {
	err := DB.Model(&model.AccountTitle{AccountTitleId: title.AccountTitleId, BookId: title.BookId}).Select("*").Updates(title).Error

	if err != nil {
		fmt.Println("Account Title could not updated: ", err)
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
	err := DB.Where(&model.BookAuthorization{UserId: user.UserId, BookId: book.BookId}).First(&bookAuthorization).Error

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
			BookId:     newBook.BookId,
			Name:       accountTitle.Name,
			Amount:     accountTitle.Amount, // 繰越
			AmountBase: accountTitle.Amount, // 繰越
			Type:       accountTitle.Type,
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

func CreateTransaction(transaction *model.Transaction) error {
	tx := DB.Begin()
	err := tx.Create(transaction).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Transaction Create Error: ", err)
		return err
	}

	for _, subTransaction := range *&transaction.SubTransactions {
		var accountTitle model.AccountTitle
		err = tx.Where(&model.AccountTitle{AccountTitleId: subTransaction.AccountTitleId}).First(&accountTitle).Error
		if err != nil {
			tx.Rollback()
			fmt.Println("Account Title not found: ", err)
			return err
		}

		if subTransaction.IsDebit {
			if accountTitle.Type%2 == 0 {
				accountTitle.Amount += subTransaction.Amount
			} else {
				accountTitle.Amount -= subTransaction.Amount
			}

		} else {
			if accountTitle.Type%2 == 0 {
				accountTitle.Amount -= subTransaction.Amount
			} else {
				accountTitle.Amount += subTransaction.Amount
			}
		}

		err = tx.Save(&accountTitle).Error
		if err != nil {
			fmt.Println("Account Title Update Error: ", err)
			tx.Rollback()
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

func UpdateTransaction(transaction *model.Transaction) error {
	var prevTransaction model.Transaction
	tx := DB.Begin()
	err := tx.Preload("SubTransactions").Where(&model.Transaction{BookId: transaction.BookId, TransactionId: transaction.TransactionId}).First(&prevTransaction).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Transaction not found: ", err)
		return err
	}

	for _, subTransaction := range prevTransaction.SubTransactions {
		var accountTitle model.AccountTitle
		err = tx.Where(&model.AccountTitle{AccountTitleId: subTransaction.AccountTitleId}).First(&accountTitle).Error
		if err != nil {
			tx.Rollback()
			fmt.Println("Account title not found: ", err)
			return err
		}

		if subTransaction.IsDebit {
			if accountTitle.Type%2 == 0 {
				accountTitle.Amount -= subTransaction.Amount
			} else {
				accountTitle.Amount += subTransaction.Amount
			}
		} else {
			if accountTitle.Type%2 == 0 {
				accountTitle.Amount += subTransaction.Amount
			} else {
				accountTitle.Amount -= subTransaction.Amount
			}
		}

		err = tx.Save(&accountTitle).Error
		if err != nil {
			fmt.Println("Account Title Update Error: ", err)
			tx.Rollback()
			return err
		}
	}

	for _, subTransaction := range transaction.SubTransactions {
		err = tx.Updates(&subTransaction).Error
		if err != nil {
			fmt.Println("Sub Transaction Update Error: ", err)
			tx.Rollback()
			return err
		}
	}

	var newTransaction model.Transaction
	err = tx.Preload("SubTransactions").Where(&model.Transaction{BookId: transaction.BookId, TransactionId: transaction.TransactionId}).First(&newTransaction).Error
	if err != nil {
		fmt.Println("Get New Transaction Update Error: ", err)
		tx.Rollback()
		return err
	}

	err = tx.Updates(transaction).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Transaction Update Error: ", err)
		return err
	}

	for _, subTransaction := range newTransaction.SubTransactions {
		var accountTitle model.AccountTitle
		err = tx.Where(&model.AccountTitle{AccountTitleId: subTransaction.AccountTitleId}).First(&accountTitle).Error
		if err != nil {
			tx.Rollback()
			fmt.Println("Account title not found: ", err)
			return err
		}

		if subTransaction.IsDebit {
			if accountTitle.Type%2 == 0 {
				accountTitle.Amount += subTransaction.Amount
			} else {
				accountTitle.Amount -= subTransaction.Amount
			}
		} else {
			if accountTitle.Type%2 == 0 {
				accountTitle.Amount -= subTransaction.Amount
			} else {
				accountTitle.Amount += subTransaction.Amount
			}
		}

		err = tx.Save(&accountTitle).Error
		if err != nil {
			fmt.Println("Account Title Update Error: ", err)
			tx.Rollback()
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

func DeleteTransaction(book *model.Book, transactionId uint64) error {
	var transaction model.Transaction

	tx := DB.Begin()
	err := tx.Preload("SubTransactions").Where(&model.Transaction{BookId: book.BookId, TransactionId: transactionId}).First(&transaction).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Transaction not found: ", err)
		return err
	}

	for _, subTransaction := range transaction.SubTransactions {
		var accountTitle model.AccountTitle
		err = tx.Where(&model.AccountTitle{AccountTitleId: subTransaction.AccountTitleId}).First(&accountTitle).Error
		if err != nil {
			tx.Rollback()
			fmt.Println("Account title not found: ", err)
			return err
		}

		if subTransaction.IsDebit {
			if accountTitle.Type%2 == 0 {
				accountTitle.Amount -= subTransaction.Amount
			} else {
				accountTitle.Amount += subTransaction.Amount
			}
		} else {
			if accountTitle.Type%2 == 0 {
				accountTitle.Amount += subTransaction.Amount
			} else {
				accountTitle.Amount -= subTransaction.Amount
			}
		}

		err = tx.Save(&accountTitle).Error
		if err != nil {
			fmt.Println("Account Title Update Error: ", err)
			tx.Rollback()
			return err
		}
	}

	err = tx.Delete(&(transaction.SubTransactions)).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Sub Transaction Delete Error: ", err)
		return err
	}

	err = tx.Delete(&transaction).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Transaction Delete Error: ", err)
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
	err := DB.Preload("SubTransactions").Preload("SubTransactions.AccountTitle").Where(&model.Transaction{TransactionId: transactionId, BookId: *&book.BookId}).First(&transaction).Error

	if err != nil {
		fmt.Println("No Transaction")
		return model.Transaction{}, err
	}

	return transaction, nil
}

func GetTransactions(book *model.Book, dataPerPage int, page int) (*[]model.Transaction, error) {
	var transactions []model.Transaction

	q := DB.Preload("SubTransactions", func(db *gorm.DB) *gorm.DB { return db.Order("sub_transactions.is_debit DESC") }).Preload("SubTransactions.AccountTitle").Where(&model.Transaction{BookId: *&book.BookId}).Order("occurred_at DESC, created_at DESC")
	err := q.Offset(dataPerPage * page).Limit(dataPerPage).Find(&transactions).Error

	if err != nil {
		fmt.Println("No Transactions", err)
		return nil, err
	}

	return &transactions, nil
}

func GetSubTransactionsFromAccountTitle(book *model.Book, accountTitleId uint64, dataPerPage int, page int) (*[]model.SubTransaction, error) {
	var subTransactions []model.SubTransaction

	err := DB.Preload("AccountTitle").Preload("Transaction", func(db *gorm.DB) *gorm.DB {
		return db.Order("transactions.occurred_at DESC, transactions.created_at DESC")
	}).Where(&model.SubTransaction{BookId: *&book.BookId, AccountTitleId: accountTitleId}).Offset(dataPerPage * page).Limit(dataPerPage).Find(&subTransactions).Error

	if err != nil {
		fmt.Println("No Sub Transactions", err)
		return nil, err
	}

	return &subTransactions, nil
}
