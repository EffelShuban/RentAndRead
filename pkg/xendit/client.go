package xendit

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/invoice"
)

type Client struct {
	api *xendit.APIClient
}

func NewClient() *Client {
	apiKey := os.Getenv("XENDIT_API_KEY")
	if apiKey == "" {
		log.Fatalf("FATAL: XENDIT_API_KEY environment variable not set")
	}
	log.Printf("Xendit Client initialized. Key length: %d", len(apiKey))
	return &Client{
		api: xendit.NewClient(apiKey),
	}
}

func (c *Client) CreateBookInvoice(ctx context.Context, externalID string, amount int, bookTitle string) (*invoice.Invoice, error) {
	desc := fmt.Sprintf("Rental for book: %s", bookTitle)
	curr := "IDR"

	data := invoice.CreateInvoiceRequest{
		ExternalId:  externalID,
		Amount:      float64(amount),
		Description: &desc,
		Currency:    &curr,
	}

	resp, httpResp, err := c.api.InvoiceApi.CreateInvoice(ctx).CreateInvoiceRequest(data).Execute()
	if err != nil {
		if httpResp != nil {
			return nil, fmt.Errorf("xendit api error (status %s): %w", httpResp.Status, err)
		}
		return nil, fmt.Errorf("xendit sdk error (no response): %v", err)
	}
	return resp, nil
}

func (c *Client) CreateTopUpInvoice(ctx context.Context, externalID string, amount int) (*invoice.Invoice, error) {
	desc := "Wallet Top-Up"
	curr := "IDR"

	data := invoice.CreateInvoiceRequest{
		ExternalId:  externalID,
		Amount:      float64(amount),
		Description: &desc,
		Currency:    &curr,
	}

	resp, httpResp, err := c.api.InvoiceApi.CreateInvoice(ctx).CreateInvoiceRequest(data).Execute()
	if err != nil {
		if httpResp != nil {
			return nil, fmt.Errorf("xendit api error (status %s): %w", httpResp.Status, err)
		}
		return nil, fmt.Errorf("xendit sdk error (no response): %v", err)
	}
	return resp, nil
}