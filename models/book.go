package model

import (
	"time"
)

type Book struct {
	BookId             string              `gorm:"default:uuid_generate_v4();primaryKey;not null;unique" json:"book_id"`
	Name               string              `gorm:"not null" json:"name"`
	Year               uint                `gorm:"not null" json:"year"`
	BookAuthorizations []BookAuthorization `gorm:"foreignKey:BookId;references:BookId;constraint:OnDelete:CASCADE;" json:"-"`
	AccountTitles      []AccountTitle      `gorm:"foreignKey:BookId;references:BookId;constraint:OnDelete:CASCADE;" json:"-"`
	Transactions       []Transaction       `gorm:"foreignKey:BookId;references:BookId;constraint:OnDelete:CASCADE;" json:"-"`
	SubTransactions    []SubTransaction    `gorm:"foreignKey:BookId;references:BookId;constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt          time.Time           `gorm:"index" json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

type BookAuthorization struct {
	BookId    string    `gorm:"primaryKey;not null"`
	Book      *Book     `gorm:"foreignKey:BookId" json:"account_title"`
	UserId    string    `gorm:"primaryKey;not null"`
	User      *User     `gorm:"foreignKey:UserId" json:"user"`
	Authority string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}

type AccountTitle struct {
	AccountTitleId  uint64           `gorm:"primaryKey;not null;autoIncrement" json:"title_id"`
	BookId          string           `gorm:"primaryKey;not null" json:"book_id"`
	Name            string           `gorm:"not null" json:"name"`
	Amount          int64            `gorm:"not null" json:"amount"`
	AmountBase      int64            `gorm:"not null;default:0" json:"amount_base"`
	SubTransactions []SubTransaction `gorm:"foreignKey:AccountTitleId,BookId;references:AccountTitleId,BookId;constraint:OnDelete:CASCADE;" json:"sub_transactions"`
	Type            uint             `gorm:"not null" json:"type"`
	CreatedAt       time.Time        `gorm:"index" json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type Transaction struct {
	TransactionId   uint64           `gorm:"index;primaryKey;not null;autoIncrement" json:"transaction_id"`
	BookId          string           `gorm:"primaryKey;not null" json:"book_id"`
	Description     string           `gorm:"not null" json:"description"`
	SubTransactions []SubTransaction `gorm:"foreignKey:TransactionId,BookId;references:TransactionId,BookId;constraint:OnDelete:CASCADE;" json:"sub_transactions"`
	OccurredAt      time.Time        `gorm:"index" json:"occurred_at"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type SubTransaction struct {
	SubTransactionId uint64        `gorm:"primaryKey;not null;autoIncrement" json:"sub_transaction_id"`
	BookId           string        `gorm:"primaryKey;not null" json:"-"`
	TransactionId    uint64        `gorm:"primaryKey;not null" json:"-"`
	Transaction      *Transaction  `gorm:"foreignKey:TransactionId,BookId" json:"transaction"`
	IsDebit          bool          `gorm:"not null" json:"is_debit"`
	AccountTitleId   uint64        `gorm:"not null" json:"account_title_id"`
	AccountTitle     *AccountTitle `gorm:"foreignKey:AccountTitleId,BookId" json:"account_title"`
	Amount           int64         `gorm:"not null" json:"amount"`
	CreatedAt        time.Time     `gorm:"index" json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}
