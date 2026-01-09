package repository

import (
	"context"
	"rent-and-read/internal/models"
	"time"

	"gorm.io/gorm"
)

type BookRental struct {
	models.Book
	RentalStatus string
	RentalID     uint
}

type BookRepository interface{
	FindByUserID(ctx context.Context, userID uint)([]BookRental, error)
	FindActiveByUserID(ctx context.Context, userID uint)([]BookRental, error)
	Find(ctx context.Context)([]models.Book, error)
	FindByID(ctx context.Context, bookID uint)(models.Book, error)
}

type psqlBookRepository struct{
	db *gorm.DB
}

func NewPsqlBookRepository(db *gorm.DB) BookRepository{
	return &psqlBookRepository{
		db: db,
	}
}

func (r *psqlBookRepository) Find(ctx context.Context)([]models.Book, error){
	var books []models.Book
	err := r.db.WithContext(ctx).Preload("Genres").Find(&books).Error
	if err != nil{
		return nil, err
	}
	return books, nil
}

func (r *psqlBookRepository) FindByUserID(ctx context.Context, userID uint)([]BookRental, error){
	var books []BookRental
	err := r.db.WithContext(ctx).
		Table("books").
		Select("books.*, rentals.status as rental_status, rentals.id as rental_id").
		Joins("JOIN rentals ON rentals.book_id = books.id").
		Where("rentals.user_id = ?", userID).
		Scan(&books).Error
	if err != nil{
		return nil, err
	}

	for i := range books {
		err := r.db.WithContext(ctx).Model(&models.Book{ID: books[i].ID}).Association("Genres").Find(&books[i].Genres)
		if err != nil {
			return nil, err
		}
	}
	return books, nil
}

func (r *psqlBookRepository) FindActiveByUserID(ctx context.Context, userID uint)([]BookRental, error){
	var books []BookRental
	err := r.db.WithContext(ctx).
		Table("books").
		Select("books.*, rentals.status as rental_status, rentals.id as rental_id").
		Joins("JOIN rentals ON rentals.book_id = books.id").
		Where("rentals.user_id = ?", userID).
		Where("rentals.status = ?", "ACTIVE").
		Where("rentals.end_date >= ?", time.Now()).
		Scan(&books).Error
	if err != nil{
		return nil, err
	}

	for i := range books {
		err := r.db.WithContext(ctx).Model(&models.Book{ID: books[i].ID}).Association("Genres").Find(&books[i].Genres)
		if err != nil {
			return nil, err
		}
	}
	return books, nil
}

func (r *psqlBookRepository) FindByID(ctx context.Context, bookID uint)(models.Book, error){
	var book models.Book
	err := r.db.WithContext(ctx).Preload("Genres").First(&book, bookID).Error
	if err != nil {
		return models.Book{}, err
	}
	return book, nil
}