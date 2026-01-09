package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"size:255;not null"`
	Email     string         `gorm:"size:255;uniqueIndex:users_email_key;not null"`
	Password  string         `gorm:"not null"` // Hashed
	Rentals   []Rental       `gorm:"foreignKey:UserID"`
	Balance int `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TopUp struct {
    gorm.Model
    UserID uint   `gorm:"not null"`
    Amount int    `gorm:"not null"`
    Status string `gorm:"size:20;default:'PENDING'"` // PENDING, SUCCESS, FAILED
}

type Genre struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"uniqueIndex:genres_name_key;not null"`
}

type Book struct {
	ID             uint           `gorm:"primaryKey"`
	ISBN           string         `gorm:"size:20;uniqueIndex:books_isbn_key;not null"`
	Title          string         `gorm:"size:255;not null"`
	Author         string         `gorm:"size:255"`
	DailyRentalFee int            `gorm:"not null"`
	TotalStock     int            `gorm:"not null;default:0"`
	AvailableStock int            `gorm:"not null;default:0"`
	Genres []Genre `gorm:"many2many:book_genres;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Rental struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null"`
	BookID    uint           `gorm:"not null"`
	StartDate time.Time      `gorm:"not null"`
	EndDate   time.Time      `gorm:"not null"`
	Status    string         `gorm:"size:20;default:'PENDING'"` // PENDING, ACTIVE, RETURNED
	NotifiedAt *time.Time
	Payment   Payment        `gorm:"foreignKey:RentalID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Payment struct {
    ID        uint           `gorm:"primaryKey"`
    RentalID         *uint   `gorm:"uniqueIndex:payments_rental_id_key"` // Nullable
    TopUpID          *uint   `gorm:"uniqueIndex:payments_topup_id_key"` // Nullable
    XenditExternalID string  `gorm:"size:255;uniqueIndex:payments_xendit_external_id_key;not null"`
    XenditInvoiceURL string  `gorm:"size:500"`
    Status           string  `gorm:"size:20;default:'UNPAID'"` // UNPAID, PAID, and SETTLED
    Amount           int     `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}