package service

import (
	"context"
	"fmt"
	"rent-and-read/internal/models"
	"rent-and-read/internal/repository"
	"time"

	"github.com/xendit/xendit-go/v6/invoice"
)

type WalletService interface {
	TopUp(ctx context.Context, userID uint, amount int) (string, string, error)
}

type XenditWalletGateway interface {
	CreateTopUpInvoice(ctx context.Context, externalID string, amount int) (*invoice.Invoice, error)
}

type walletService struct {
	xenditClient XenditWalletGateway
	userRepo     repository.UserRepository
}

func NewWalletService(xenditClient XenditWalletGateway, userRepo repository.UserRepository) WalletService {
	return &walletService{
		xenditClient: xenditClient,
		userRepo:     userRepo,
	}
}

func (s *walletService) TopUp(ctx context.Context, userID uint, amount int) (string, string, error) {
	externalID := fmt.Sprintf("TOPUP-%d-%d", userID, time.Now().Unix())

	invoice, err := s.xenditClient.CreateTopUpInvoice(ctx, externalID, amount)
	if err != nil {
		return "", "", fmt.Errorf("xendit error: %+v", err)
	}

	topUp := &models.TopUp{
		UserID: userID,
		Amount: amount,
		Status: "PENDING",
	}

	payment := &models.Payment{
		Amount:           amount,
		XenditExternalID: externalID,
		XenditInvoiceURL: invoice.GetInvoiceUrl(),
		Status:           "UNPAID",
	}

	err = s.userRepo.CreateTopUp(ctx, topUp, payment)
	if err != nil {
		return "", "", fmt.Errorf("database error: %v", err)
	}

	return invoice.GetInvoiceUrl(), externalID, nil
}