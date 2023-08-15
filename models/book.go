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
	UserId    string    `gorm:"primaryKey;not null"`
	Authority string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}

type AccountTitle struct {
	AccountTitleId  uint64           `gorm:"primaryKey;not null;autoIncrement" json:"title_id"`
	BookId          string           `gorm:"primaryKey;not null" json:"book_id"`
	Name            string           `gorm:"not null" json:"name"`
	Amount          int64            `gorm:"not null" json:"amount"`
	Type            uint             `gorm:"not null" json:"type"`
	Transactions    []Transaction    `gorm:"foreignKey:DebitId,CreditId;references:AccountTitleId,AccountTitleId;constraint:OnDelete:CASCADE;" json:"-"`
	SubTransactions []SubTransaction `gorm:"foreignKey:DebitId,CreditId;references:AccountTitleId,AccountTitleId;constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt       time.Time        `gorm:"index" json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type Transaction struct {
	TransactionId   uint64           `gorm:"primaryKey;not null;autoIncrement"`
	BookId          string           `gorm:"primaryKey;not null"`
	DebitId         uint64           `gorm:"not null"`
	CreditId        uint64           `gorm:"not null"`
	Description     string           `gorm:"not null"`
	Amount          int64            `gorm:"not null"`
	SubTransactions []SubTransaction `gorm:"foreignKey:TransactionId;references:TransactionId;constraint:OnDelete:CASCADE;"`
	CreatedAt       time.Time        `gorm:"index"`
	UpdatedAt       time.Time
}

type SubTransaction struct {
	SubTransactionId uint64    `gorm:"primaryKey;not null;autoIncrement"`
	BookId           string    `gorm:"primaryKey;not null"`
	TransactionId    uint64    `gorm:"primaryKey;not null"`
	DebitId          uint64    `gorm:"not null"`
	CreditId         uint64    `gorm:"not null"`
	Amount           int64     `gorm:"not null"`
	CreatedAt        time.Time `gorm:"index"`
	UpdatedAt        time.Time
}
