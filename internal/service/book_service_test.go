package service

import (
	"context"
	"errors"
	"rent-and-read/internal/models"
	"rent-and-read/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookService_Find(t *testing.T) {
	mockRepo := new(MockBookRepo)
	service := NewBookService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		expectedBooks := []models.Book{
			{ID: 1, Title: "Book 1", Genres: []models.Genre{{Name: "Fiction"}}},
			{ID: 2, Title: "Book 2"},
		}
		mockRepo.On("Find", ctx).Return(expectedBooks, nil).Once()

		result, err := service.Find(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Book 1", result[0].Title)
		assert.Equal(t, "Fiction", result[0].Genres[0])
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("Find", ctx).Return([]models.Book(nil), errors.New("db error")).Once()

		result, err := service.Find(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookService_FindByUserID(t *testing.T) {
	mockRepo := new(MockBookRepo)
	service := NewBookService(mockRepo)
	ctx := context.Background()
	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		expectedBooks := []repository.BookRental{
			{
				Book: models.Book{ID: 1, Title: "Rented Book", Genres: []models.Genre{{Name: "Sci-Fi"}}},
				RentalStatus: "ACTIVE",
				RentalID:     123,
			},
		}
		mockRepo.On("FindByUserID", ctx, userID).Return(expectedBooks, nil).Once()

		result, err := service.FindByUserID(ctx, userID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Rented Book", result[0].Title)
		assert.Equal(t, "ACTIVE", result[0].Status)
		assert.Equal(t, uint(123), result[0].RentalID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("FindByUserID", ctx, userID).Return([]repository.BookRental(nil), errors.New("db error")).Once()

		result, err := service.FindByUserID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}