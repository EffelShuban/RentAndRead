package handler

import (
	"fmt"
	"net/http"
	"rent-and-read/internal/dto"
	"rent-and-read/internal/service"
	"rent-and-read/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RentalHandler struct {
	serv service.RentalService
}

func NewRentalHandler(serv service.RentalService) *RentalHandler {
	return &RentalHandler{serv: serv}
}

// CreateRental godoc
// @Summary Rent a book
// @Description Create a new rental transaction using Wallet or Xendit
// @Tags Rentals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RentalRequest true "Rental Request"
// @Success 201 {object} dto.RentalResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/rentals [post]
func (h *RentalHandler) CreateRental(c echo.Context) error {
	var req dto.RentalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}
	
	if req.BookID == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "book_id is required"})
	}
	if req.Days < 1 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "days must be at least 1"})
	}

	userID := utils.GetUserID(c)
	if userID == 0{
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid token format"})
	}

	resp, err := h.serv.Create(c.Request().Context(), userID, req.BookID, req.Days)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message":     resp.Message,
		"payment_url": resp.PaymentURL,
		"amount":      resp.Amount,
		"paid_via":    resp.PaidVia,
		"external_id": resp.XenditExternalID,
	})
}

// ReturnRental godoc
// @Summary Return a rented book
// @Description Return a book and update stock
// @Tags Rentals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Rental ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/rentals/{id}/return [post]
func (h *RentalHandler) ReturnRental(c echo.Context) error {
	userID := utils.GetUserID(c)
	rentalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid rental id"})
	}

	lateFee, err := h.serv.Return(c.Request().Context(), userID, uint(rentalID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	msg := "book returned successfully"
	if lateFee > 0 {
		msg = fmt.Sprintf("book returned successfully with late fee of %d deducted", lateFee)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": msg, "late_fee": lateFee})
}