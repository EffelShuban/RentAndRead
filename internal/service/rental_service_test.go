package service

import (
	"context"
	"errors"
	"rent-and-read/internal/models"
	"rent-and-read/internal/repository"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xendit/xendit-go/v6/invoice"
)

type MockBookRepo struct {
	mock.Mock
}

func (m *MockBookRepo) FindByID(ctx context.Context, bookID uint) (models.Book, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(models.Book), args.Error(1)
}
func (m *MockBookRepo) Find(ctx context.Context) ([]models.Book, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Book), args.Error(1)
}
func (m *MockBookRepo) FindByUserID(ctx context.Context, userID uint) ([]repository.BookRental, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]repository.BookRental), args.Error(1)
}

func (m *MockBookRepo) FindActiveByUserID(ctx context.Context, userID uint) ([]repository.BookRental, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]repository.BookRental), args.Error(1)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByID(ctx context.Context, userID uint) (models.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.User), args.Error(1)
}
func (m *MockUserRepo) Create(ctx context.Context, u models.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}
func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}
func (m *MockUserRepo) CreateTopUp(ctx context.Context, topUp *models.TopUp, payment *models.Payment) error {
	args := m.Called(ctx, topUp, payment)
	return args.Error(0)
}
func (m *MockUserRepo) HandlePaymentSuccess(ctx context.Context, externalID string) error {
	args := m.Called(ctx, externalID)
	return args.Error(0)
}

type MockRentalRepo struct {
	mock.Mock
}

func (m *MockRentalRepo) CreateWithWallet(ctx context.Context, userID uint, bookID uint, totalFee int, rental *models.Rental) error {
	args := m.Called(ctx, userID, bookID, totalFee, rental)
	return args.Error(0)
}
func (m *MockRentalRepo) Create(ctx context.Context, rental *models.Rental, payment *models.Payment) error {
	args := m.Called(ctx, rental, payment)
	return args.Error(0)
}
func (m *MockRentalRepo) Return(ctx context.Context, userID uint, rentalID uint) (int, error) {
	args := m.Called(ctx, userID, rentalID)
	return args.Int(0), args.Error(1)
}

func (m *MockRentalRepo) FindRentalsDueIn(ctx context.Context, days int) ([]models.Rental, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]models.Rental), args.Error(1)
}

func (m *MockRentalRepo) MarkAsNotified(ctx context.Context, rentalID uint) error {
	args := m.Called(ctx, rentalID)
	return args.Error(0)
}

type MockXenditClient struct {
	mock.Mock
}

func (m *MockXenditClient) CreateBookInvoice(ctx context.Context, externalID string, amount int, bookTitle string) (*invoice.Invoice, error) {
	args := m.Called(ctx, externalID, amount, bookTitle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*invoice.Invoice), args.Error(1)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendRentalDueReminder(toEmail, name, bookTitle string, dueDate time.Time) error {
	args := m.Called(toEmail, name, bookTitle, dueDate)
	return args.Error(0)
}

func TestCreateRental_Wallet_Success(t *testing.T) {
	mockBookRepo := new(MockBookRepo)
	mockUserRepo := new(MockUserRepo)
	mockRentalRepo := new(MockRentalRepo)
	mockXenditClient := new(MockXenditClient)
	mockEmailService := new(MockEmailService)

	service := NewRentalService(mockRentalRepo, mockBookRepo, mockUserRepo, mockXenditClient, mockEmailService)

	ctx := context.Background()
	userID := uint(1)
	bookID := uint(10)
	days := 3
	dailyFee := 5000
	totalFee := dailyFee * days

	mockBookRepo.On("FindByID", ctx, bookID).Return(models.Book{
		ID:             bookID,
		Title:          "Test Book",
		DailyRentalFee: dailyFee,
		AvailableStock: 5,
	}, nil)

	mockUserRepo.On("FindByID", ctx, userID).Return(models.User{
		ID:      userID,
		Balance: 20000,
	}, nil)

	mockRentalRepo.On("CreateWithWallet", ctx, userID, bookID, totalFee, mock.AnythingOfType("*models.Rental")).Return(nil)

	resp, err := service.Create(ctx, userID, bookID, days)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "WALLET", resp.PaidVia)
	assert.Equal(t, totalFee, resp.Amount)

	mockBookRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockRentalRepo.AssertExpectations(t)
}

func TestCreateRental_BookNotFound(t *testing.T) {
	mockBookRepo := new(MockBookRepo)
	service := NewRentalService(nil, mockBookRepo, nil, nil, nil)

	ctx := context.Background()
	mockBookRepo.On("FindByID", ctx, uint(1)).Return(models.Book{}, errors.New("book not found"))

	resp, err := service.Create(ctx, 1, 1, 3)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "book not found", err.Error())
	mockBookRepo.AssertExpectations(t)
}

func TestCreateRental_UserNotFound(t *testing.T) {
	mockBookRepo := new(MockBookRepo)
	mockUserRepo := new(MockUserRepo)
	service := NewRentalService(nil, mockBookRepo, mockUserRepo, nil, nil)

	ctx := context.Background()
	mockBookRepo.On("FindByID", ctx, uint(10)).Return(models.Book{DailyRentalFee: 5000}, nil)
	mockUserRepo.On("FindByID", ctx, uint(1)).Return(models.User{}, errors.New("user not found"))

	resp, err := service.Create(ctx, 1, 10, 3)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "user not found", err.Error())
	mockBookRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateRental_Wallet_RepoError(t *testing.T) {
	mockBookRepo := new(MockBookRepo)
	mockUserRepo := new(MockUserRepo)
	mockRentalRepo := new(MockRentalRepo)
	service := NewRentalService(mockRentalRepo, mockBookRepo, mockUserRepo, nil, nil)

	ctx := context.Background()
	mockBookRepo.On("FindByID", ctx, uint(10)).Return(models.Book{DailyRentalFee: 5000}, nil)
	mockUserRepo.On("FindByID", ctx, uint(1)).Return(models.User{Balance: 20000}, nil)
	mockRentalRepo.On("CreateWithWallet", ctx, uint(1), uint(10), 15000, mock.Anything).Return(errors.New("db error"))

	resp, err := service.Create(ctx, 1, 10, 3)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "db error", err.Error())
	mockRentalRepo.AssertExpectations(t)
}

func TestCreateRental_Xendit_Success(t *testing.T) {
	mockBookRepo := new(MockBookRepo)
	mockUserRepo := new(MockUserRepo)
	mockRentalRepo := new(MockRentalRepo)
	mockXenditClient := new(MockXenditClient)
	mockEmailService := new(MockEmailService)
	service := NewRentalService(mockRentalRepo, mockBookRepo, mockUserRepo, mockXenditClient, mockEmailService)

	ctx := context.Background()
	bookID := uint(10)
	userID := uint(1)
	totalFee := 15000

	mockBookRepo.On("FindByID", ctx, bookID).Return(models.Book{Title: "Go Book", DailyRentalFee: 5000}, nil)
	mockUserRepo.On("FindByID", ctx, userID).Return(models.User{Balance: 5000}, nil) // Insufficient balance

	mockXenditClient.On("CreateBookInvoice", ctx, mock.MatchedBy(func(id string) bool {
		return strings.HasPrefix(id, "RENT-1-")
	}), totalFee, "Go Book").Return(&invoice.Invoice{}, nil)

	mockRentalRepo.On("Create", ctx, mock.AnythingOfType("*models.Rental"), mock.AnythingOfType("*models.Payment")).Return(nil)

	resp, err := service.Create(ctx, userID, bookID, 3)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "XENDIT", resp.PaidVia)
	assert.Equal(t, totalFee, resp.Amount)
	assert.Contains(t, resp.XenditExternalID, "RENT-1-")

	mockXenditClient.AssertExpectations(t)
	mockRentalRepo.AssertExpectations(t)
}

func TestCreateRental_Xendit_InvoiceError(t *testing.T) {
	mockBookRepo := new(MockBookRepo)
	mockUserRepo := new(MockUserRepo)
	mockXenditClient := new(MockXenditClient)
	service := NewRentalService(nil, mockBookRepo, mockUserRepo, mockXenditClient, nil)

	ctx := context.Background()
	mockBookRepo.On("FindByID", ctx, uint(10)).Return(models.Book{Title: "Go Book", DailyRentalFee: 5000}, nil)
	mockUserRepo.On("FindByID", ctx, uint(1)).Return(models.User{Balance: 5000}, nil)

	mockXenditClient.On("CreateBookInvoice", ctx, mock.Anything, 15000, "Go Book").Return(nil, errors.New("xendit error"))

	resp, err := service.Create(ctx, 1, 10, 3)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "xendit error", err.Error())
	mockXenditClient.AssertExpectations(t)
}

func TestCreateRental_Xendit_RepoError(t *testing.T) {
	mockBookRepo := new(MockBookRepo)
	mockUserRepo := new(MockUserRepo)
	mockRentalRepo := new(MockRentalRepo)
	mockXenditClient := new(MockXenditClient)
	mockEmailService := new(MockEmailService)
	service := NewRentalService(mockRentalRepo, mockBookRepo, mockUserRepo, mockXenditClient, mockEmailService)

	ctx := context.Background()
	mockBookRepo.On("FindByID", ctx, uint(10)).Return(models.Book{Title: "Go Book", DailyRentalFee: 5000}, nil)
	mockUserRepo.On("FindByID", ctx, uint(1)).Return(models.User{Balance: 5000}, nil)

	mockXenditClient.On("CreateBookInvoice", ctx, mock.Anything, 15000, "Go Book").Return(&invoice.Invoice{}, nil)
	mockRentalRepo.On("Create", ctx, mock.Anything, mock.Anything).Return(errors.New("db error"))

	resp, err := service.Create(ctx, 1, 10, 3)

	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Invoice created", resp.Message)
	assert.Equal(t, "db error", err.Error())
	mockRentalRepo.AssertExpectations(t)
}

func TestSendDueReminders_Success(t *testing.T) {
	mockRentalRepo := new(MockRentalRepo)
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	mockEmailService := new(MockEmailService)
	service := NewRentalService(mockRentalRepo, mockBookRepo, mockUserRepo, nil, mockEmailService)

	ctx := context.Background()
	now := time.Now()
	rentals := []models.Rental{
		{ID: 1, UserID: 1, BookID: 10, EndDate: now.Add(24 * time.Hour)},
	}

	mockRentalRepo.On("FindRentalsDueIn", ctx, 1).Return(rentals, nil)

	mockUserRepo.On("FindByID", ctx, uint(1)).Return(models.User{ID: 1, Name: "John", Email: "john@example.com"}, nil)
	mockBookRepo.On("FindByID", ctx, uint(10)).Return(models.Book{ID: 10, Title: "Go Guide"}, nil)

	mockEmailService.On("SendRentalDueReminder", "john@example.com", "John", "Go Guide", rentals[0].EndDate).Return(nil)

	mockRentalRepo.On("MarkAsNotified", ctx, rentals[0].ID).Return(nil)

	err := service.SendDueReminders(ctx)

	assert.NoError(t, err)
	mockRentalRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}