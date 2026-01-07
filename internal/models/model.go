package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"size:255;not null"`
	Email     string         `gorm:"size:255;uniqueIndex;not null"`
	Password  string         `gorm:"not null"` // Hashed
	Rentals   []Rental       `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Book model for inventory and stock management
type Book struct {
	ID             uint           `gorm:"primaryKey"`
	ISBN           string         `gorm:"size:20;uniqueIndex;not null"`
	Title          string         `gorm:"size:255;not null"`
	Author         string         `gorm:"size:255"`
	DailyRentalFee int            `gorm:"not null"` // Price in smallest unit (e.g., IDR)
	TotalStock     int            `gorm:"not null;default:0"`
	AvailableStock int            `gorm:"not null;default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// Rental model to track book borrowing sessions
type Rental struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null"`
	BookID    uint           `gorm:"not null"`
	StartDate time.Time      `gorm:"not null"`
	EndDate   time.Time      `gorm:"not null"`
	Status    string         `gorm:"size:20;default:'PENDING'"` // PENDING, ACTIVE, RETURNED
	Payment   Payment        `gorm:"foreignKey:RentalID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Payment model to integrate with Xendit Invoice
type Payment struct {
	ID               uint           `gorm:"primaryKey"`
	RentalID         uint           `gorm:"uniqueIndex;not null"`
	XenditExternalID string         `gorm:"size:255;uniqueIndex;not null"`
	XenditInvoiceURL string         `gorm:"size:500"`
	Status           string         `gorm:"size:20;default:'UNPAID'"` // UNPAID, PAID, EXPIRED
	Amount           int            `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}