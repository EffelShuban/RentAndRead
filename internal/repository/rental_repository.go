package repository

import (
	"context"
	"errors"
	"fmt"
	"rent-and-read/internal/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RentalRepository interface {
	Create(ctx context.Context, rental *models.Rental, payment *models.Payment) error 
	CreateWithWallet(ctx context.Context, userID uint, bookID uint, totalFee int, rental *models.Rental) error
	Return(ctx context.Context, userID uint, rentalID uint) (int, error)
	FindRentalsDueIn(ctx context.Context, days int) ([]models.Rental, error)
	MarkAsNotified(ctx context.Context, rentalID uint) error
}

type psqlRentalRepository struct {
	db *gorm.DB
}

func NewPsqlRentalRepository(db *gorm.DB) RentalRepository {
	return &psqlRentalRepository{db: db}
}

func (r *psqlRentalRepository) Create(ctx context.Context, rental *models.Rental, payment *models.Payment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var book models.Book
		
		if err := tx.First(&book, rental.BookID).Error; err != nil {
			return err
		}

		if book.AvailableStock <= 0 {
			return models.ErrBookOutOfStock
		}

		if err := tx.Create(rental).Error; err != nil {
			return err
		}

		payment.RentalID = &rental.ID
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *psqlRentalRepository) CreateWithWallet(ctx context.Context, userID uint, bookID uint, totalFee int, rental *models.Rental) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var user models.User
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
            return err
        }

        if user.Balance < totalFee {
            return models.ErrInsufficientBalance
        }

        var book models.Book
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, bookID).Error; err != nil {
            return err
        }

        if book.AvailableStock <= 0 {
            return models.ErrBookOutOfStock
        }

        if err := tx.Model(&user).Update("balance", user.Balance-totalFee).Error; err != nil {
            return err
        }

        if err := tx.Model(&book).Update("available_stock", book.AvailableStock-1).Error; err != nil {
            return err
        }

        rental.Status = "ACTIVE"
        return tx.Create(rental).Error
    })
}

func (r *psqlRentalRepository) Return(ctx context.Context, userID uint, rentalID uint) (int, error) {
	var lateFee int
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rental models.Rental
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", rentalID, userID).
			First(&rental).Error; err != nil {
			return err
		}

		if rental.Status != "ACTIVE" {
			return errors.New("rental is not active")
		}

		if time.Now().After(rental.EndDate) {
			daysLate := int(time.Since(rental.EndDate).Hours()/24) + 1
			
			var book models.Book
			if err := tx.First(&book, rental.BookID).Error; err != nil {
				return err
			}
			
			lateFee = daysLate * book.DailyRentalFee

			var user models.User
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
				return err
			}

			if user.Balance < lateFee {
				return fmt.Errorf("insufficient balance to pay late fee of %d", lateFee)
			}

			if err := tx.Model(&user).Update("balance", user.Balance - lateFee).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&rental).Update("status", "RETURNED").Error; err != nil {
			return err
		}

		return tx.Model(&models.Book{ID: rental.BookID}).
			Update("available_stock", gorm.Expr("available_stock + ?", 1)).Error
	})
	return lateFee, err
}

func (r *psqlRentalRepository) FindRentalsDueIn(ctx context.Context, days int) ([]models.Rental, error) {
	var rentals []models.Rental
	target := time.Now().AddDate(0, 0, days)
	start := time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, target.Location())
	end := start.Add(24 * time.Hour)

	err := r.db.WithContext(ctx).
		Where("status = ? AND end_date >= ? AND end_date < ? AND notified_at IS NULL", "ACTIVE", start, end).
		Find(&rentals).Error
	return rentals, err
}

func (r *psqlRentalRepository) MarkAsNotified(ctx context.Context, rentalID uint) error {
	return r.db.WithContext(ctx).Model(&models.Rental{}).Where("id = ?", rentalID).Update("notified_at", time.Now()).Error
}