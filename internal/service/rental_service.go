package service

import (
	"context"
	"fmt"
	"rent-and-read/internal/dto"
	"rent-and-read/internal/models"
	"rent-and-read/internal/repository"
	"time"

	"github.com/xendit/xendit-go/v6/invoice"
)

type RentalService interface {
	Create(ctx context.Context, userID, bookID uint, days int) (*dto.RentalResponse, error)
	Return(ctx context.Context, userID, rentalID uint) (int, error)
	SendDueReminders(ctx context.Context) error
}

type XenditPaymentGateway interface {
	CreateBookInvoice(ctx context.Context, externalID string, amount int, bookTitle string) (*invoice.Invoice, error)
}

type EmailService interface {
	SendRentalDueReminder(toEmail, name, bookTitle string, dueDate time.Time) error
}

type rentalService struct {
	rentalRepo   repository.RentalRepository
	bookRepo     repository.BookRepository
	userRepo repository.UserRepository
	xenditClient XenditPaymentGateway
	emailService EmailService
}

func NewRentalService(rentalRepo repository.RentalRepository, bookRepo repository.BookRepository, userRepo repository.UserRepository, xenditClient XenditPaymentGateway, emailService EmailService) RentalService {
	return &rentalService{
		rentalRepo: rentalRepo,
		bookRepo: bookRepo,
		userRepo: userRepo,
		xenditClient: xenditClient,
		emailService: emailService,
	}
}

func (s *rentalService) Create(ctx context.Context, userID, bookID uint, days int) (*dto.RentalResponse, error) {
    book, err := s.bookRepo.FindByID(ctx, bookID)
    if err != nil {
        return nil, err
    }
    totalAmount := book.DailyRentalFee * days
    
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }

    if user.Balance >= totalAmount {
        rental := &models.Rental{
            UserID: userID, BookID: bookID, 
            StartDate: time.Now(), EndDate: time.Now().AddDate(0,0,days),
        }
        err := s.rentalRepo.CreateWithWallet(ctx, userID, bookID, totalAmount, rental)
        if err != nil { return nil, err }
        return &dto.RentalResponse{Message: "Success", PaidVia: "WALLET", Amount: totalAmount}, nil
    }

    externalID := fmt.Sprintf("RENT-%d-%d", userID, time.Now().Unix())
    invoice, err := s.xenditClient.CreateBookInvoice(ctx, externalID, totalAmount, book.Title)
    if err != nil { return nil, err }

    payment := &models.Payment{
        Amount: totalAmount, XenditExternalID: externalID, 
        XenditInvoiceURL: invoice.InvoiceUrl, Status: "UNPAID",
    }
    
    rental := &models.Rental{
        UserID: userID, 
        BookID: bookID, 
        Status: "PENDING",
        StartDate: time.Now(),
        EndDate: time.Now().AddDate(0, 0, days),
    }
    
    err = s.rentalRepo.Create(ctx, rental, payment)
    return &dto.RentalResponse{
        Message: "Invoice created", PaymentURL: invoice.InvoiceUrl, PaidVia: "XENDIT", Amount: totalAmount, XenditExternalID: externalID,
    }, err
}

func (s *rentalService) Return(ctx context.Context, userID, rentalID uint) (int, error) {
	return s.rentalRepo.Return(ctx, userID, rentalID)
}

func (s *rentalService) SendDueReminders(ctx context.Context) error {
	rentals, err := s.rentalRepo.FindRentalsDueIn(ctx, 1)
	if err != nil {
		return err
	}

	for _, r := range rentals {
		user, err := s.userRepo.FindByID(ctx, r.UserID)
		if err != nil { continue }

		book, err := s.bookRepo.FindByID(ctx, r.BookID)
		if err != nil { continue }

		if err := s.emailService.SendRentalDueReminder(user.Email, user.Name, book.Title, r.EndDate); err == nil {
			_ = s.rentalRepo.MarkAsNotified(ctx, r.ID)
		}
	}

	return nil
}