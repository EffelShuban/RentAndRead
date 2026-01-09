package dto

type BookResponse struct {
	ID             uint     `json:"id"`
	ISBN           string   `json:"isbn"`
	Title          string   `json:"title"`
	Author         string   `json:"author"`
	AvailableStock int      `json:"available_stock"`
	Genres         []string `json:"genres"`
	Status         string   `json:"status,omitempty"`
	RentalID       uint     `json:"rental_id,omitempty"`
}

type RegisterRequest struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct{
	Name string `json:"name"`
	Email string `json:"email"`
}

type RentalRequest struct {
	BookID uint `json:"book_id" validate:"required"`
	Days   int  `json:"days" validate:"required,min=1"`
}

type RentalResponse struct {
    Message          string `json:"message"`
    PaymentURL       string `json:"payment_url,omitempty"`
    XenditExternalID string `json:"external_id,omitempty"`
    Amount           int    `json:"amount"`
    PaidVia          string `json:"paid_via"` // "WALLET" or "XENDIT"
}

type TopUpRequest struct {
	Amount int `json:"amount" validate:"required,min=10000"`
}

type XenditCallbackRequest struct {
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
}