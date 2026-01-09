package service

import (
	"context"
	"rent-and-read/internal/dto"
	"rent-and-read/internal/repository"
)

type BookService interface{
	FindByUserID(ctx context.Context, userID uint) ([]dto.BookResponse, error)
	FindActiveByUserID(ctx context.Context, userID uint) ([]dto.BookResponse, error)
	Find(ctx context.Context) ([]dto.BookResponse, error)
}

type bookService struct{
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService{
	return &bookService{
		repo: repo,
	}
}

func (s *bookService) Find(ctx context.Context) ([]dto.BookResponse, error) {
	books, err := s.repo.Find(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.BookResponse, 0)
	for _, b := range books {
		genreNames := make([]string, len(b.Genres))
		for i, g := range b.Genres {
			genreNames[i] = g.Name
		}

		response = append(response, dto.BookResponse{
			ID:             b.ID,
			ISBN:           b.ISBN,
			Title:          b.Title,
			Author:         b.Author,
			AvailableStock: b.AvailableStock,
			Genres:         genreNames,
		})
	}

	return response, nil
}

func (s *bookService) FindActiveByUserID(ctx context.Context, userID uint) ([]dto.BookResponse, error){
	books, err := s.repo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.BookResponse, 0)
	for _, b := range books {
		genreNames := make([]string, len(b.Genres))
		for i, g := range b.Genres {
			genreNames[i] = g.Name
		}

		response = append(response, dto.BookResponse{
			ID:             b.ID,
			ISBN:           b.ISBN,
			Title:          b.Title,
			Author:         b.Author,
			AvailableStock: b.AvailableStock,
			Genres:         genreNames,
			Status:         b.RentalStatus,
			RentalID:       b.RentalID,
		})
	}

	return response, nil
}

func (s *bookService) FindByUserID(ctx context.Context, userID uint) ([]dto.BookResponse, error){
	books, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.BookResponse, 0)
	for _, b := range books {
		genreNames := make([]string, len(b.Genres))
		for i, g := range b.Genres {
			genreNames[i] = g.Name
		}

		response = append(response, dto.BookResponse{
			ID:             b.ID,
			ISBN:           b.ISBN,
			Title:          b.Title,
			Author:         b.Author,
			AvailableStock: b.AvailableStock,
			Genres:         genreNames,
			Status:         b.RentalStatus,
			RentalID:       b.RentalID,
		})
	}

	return response, nil
}