package repository

import (
	"context"
	"errors"
	"rent-and-read/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface{
	Create(ctx context.Context, u models.User) error
	FindByID(ctx context.Context, userID uint) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	CreateTopUp(ctx context.Context, topUp *models.TopUp, payment *models.Payment) error
	HandlePaymentSuccess(ctx context.Context, externalID string) error
}

type psqlUserRepository struct{
	db *gorm.DB
}

func NewPsqlUserRepository(db *gorm.DB) UserRepository{
	return &psqlUserRepository{
		db: db,
	}
}

func (r *psqlUserRepository) FindByEmail(ctx context.Context, email string) (models.User, error){
	var u models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, err
	}

	return u, nil
}

func (r *psqlUserRepository) FindByID(ctx context.Context, userID uint) (models.User, error){
	var u models.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&u).Error
	if err != nil{
		return models.User{}, models.ErrUserNotFound
	}

	return u, nil
}

func (r *psqlUserRepository) Create(ctx context.Context, u models.User) error{
	_, err := r.FindByEmail(ctx, u.Email)
	if err == nil{
		return models.ErrUserExist
	}

	return r.db.Create(&u).Error
}

func (r *psqlUserRepository) CreateTopUp(ctx context.Context, topUp *models.TopUp, payment *models.Payment) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(topUp).Error; err != nil { return err }
        payment.TopUpID = &topUp.ID
        return tx.Create(payment).Error
    })
}

func (r *psqlUserRepository) HandlePaymentSuccess(ctx context.Context, externalID string) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var payment models.Payment
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("xendit_external_id = ?", externalID).First(&payment).Error; err != nil {
            return err
        }

        if payment.Status == "PAID"{
            return nil
        }

        if err := tx.Model(&payment).Update("status", "PAID").Error; err != nil {
            return err
        }

        if payment.TopUpID != nil {
            var topUp models.TopUp
            if err := tx.First(&topUp, *payment.TopUpID).Error; err != nil {
                return err
            }
            if err := tx.Model(&topUp).Update("status", "SUCCESS").Error; err != nil {
                return err
            }
            
            return tx.Model(&models.User{}).Where("id = ?", topUp.UserID).
                Update("balance", gorm.Expr("balance + ?", payment.Amount)).Error

        } else if payment.RentalID != nil {
            var rental models.Rental
            if err := tx.First(&rental, *payment.RentalID).Error; err != nil {
                return err
            }

            var book models.Book
            if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, rental.BookID).Error; err != nil {
                return err
            }

            if book.AvailableStock > 0 {
                if err := tx.Model(&book).Update("available_stock", book.AvailableStock-1).Error; err != nil {
                    return err
                }
                return tx.Model(&rental).Update("status", "ACTIVE").Error
            } else {
                return tx.Model(&rental).Update("status", "FAILED").Error
            }
        }

        return nil
    })
}