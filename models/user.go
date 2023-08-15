package model

import (
	"time"
)

type User struct {
	UserId             string              `gorm:"default:uuid_generate_v4();primaryKey;not nullunique"`
	Email              string              `gorm:"not null;unique"`
	Name               string              `gorm:"not null"`
	Password           string              `gorm:"not null"`
	BookAuthorizations []BookAuthorization `gorm:"foreignKey:UserId;references:UserId;constraint:OnDelete:CASCADE;"`
	Applications       []Application       `gorm:"foreignKey:UserId;references:UserId;constraint:OnDelete:CASCADE;"`
	Permits            []Permit            `gorm:"foreignKey:UserId;references:UserId;constraint:OnDelete:CASCADE;"`
	CreatedAt          time.Time           `gorm:"index"`
	UpdatedAt          time.Time
}

type Application struct {
	ApplicationId      string    `gorm:"default:uuid_generate_v4();primaryKey;not null;unique"`
	UserId             uint64    `gorm:"not null"`
	secret             string    `gorm:"not null"`
	name               string    `gorm:"not null"`
	auth_redirect_uri  string    `gorm:"not null"`
	token_redirect_uri string    `gorm:"not null"`
	Permits            []Permit  `gorm:"foreignKey:ApplicationId;references:ApplicationId;constraint:OnDelete:CASCADE;"`
	CreatedAt          time.Time `gorm:"index"`
	UpdatedAt          time.Time
}

type Permit struct {
	UserId        string    `gorm:"primaryKey;not null"`
	ApplicationId string    `gorm:"primaryKey;not null"`
	Authority     string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"index"`
	UpdatedAt     time.Time
}
