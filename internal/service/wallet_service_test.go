package service

import (
	"context"
	"errors"
	"rent-and-read/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xendit/xendit-go/v6/invoice"
)

type MockXenditWalletGateway struct {
	mock.Mock
}

func (m *MockXenditWalletGateway) CreateTopUpInvoice(ctx context.Context, externalID string, amount int) (*invoice.Invoice, error) {
	args := m.Called(ctx, externalID, amount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*invoice.Invoice), args.Error(1)
}

func TestWalletService_TopUp(t *testing.T) {
	mockXendit := new(MockXenditWalletGateway)
	mockRepo := new(MockUserRepo)
	service := NewWalletService(mockXendit, mockRepo)
	ctx := context.Background()
	userID := uint(1)
	amount := 50000

	t.Run("Success", func(t *testing.T) {
		mockInvoice := &invoice.Invoice{
			InvoiceUrl: "https://checkout.xendit.co/web/123",
		}

		mockXendit.On("CreateTopUpInvoice", ctx, mock.AnythingOfType("string"), amount).Return(mockInvoice, nil).Once()

		mockRepo.On("CreateTopUp", ctx, mock.MatchedBy(func(t *models.TopUp) bool {
			return t.UserID == userID && t.Amount == amount
		}), mock.MatchedBy(func(p *models.Payment) bool {
			return p.Amount == amount && p.XenditInvoiceURL == mockInvoice.InvoiceUrl
		})).Return(nil).Once()

		url, extID, err := service.TopUp(ctx, userID, amount)

		assert.NoError(t, err)
		assert.Equal(t, mockInvoice.InvoiceUrl, url)
		assert.NotEmpty(t, extID)
		mockXendit.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("XenditError", func(t *testing.T) {
		mockXendit.On("CreateTopUpInvoice", ctx, mock.Anything, amount).Return(nil, errors.New("api error")).Once()

		url, extID, err := service.TopUp(ctx, userID, amount)

		assert.Error(t, err)
		assert.Empty(t, url)
		assert.Empty(t, extID)
	})
}