package handler

import (
	"errors"
	"net/http"
	"rent-and-read/internal/dto"
	"rent-and-read/internal/models"
	"rent-and-read/internal/service"
	"rent-and-read/utils"

	"github.com/labstack/echo/v4"
)

type WalletHandler struct {
	userServ   service.UserService
	walletServ service.WalletService
}

func NewWalletHandler(userServ service.UserService, walletServ service.WalletService) *WalletHandler {
	return &WalletHandler{
		userServ:   userServ,
		walletServ: walletServ,
	}
}

// TopUp godoc
// @Summary Top up wallet balance
// @Description Request a top-up for the user's wallet
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TopUpRequest true "TopUp Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/wallet/topup [post]
func (h *WalletHandler) TopUp(c echo.Context) error {
	var req dto.TopUpRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid token"})
	}

	// This calls the WalletService to get the Xendit URL
	paymentURL, externalID, err := h.walletServ.TopUp(c.Request().Context(), userID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":     "Top-up invoice created",
		"payment_url": paymentURL,
		"external_id": externalID,
	})
}

// HandleXenditCallback godoc
// @Summary Handle Xendit Payment Callback
// @Description Webhook endpoint for Xendit payment status updates
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param request body dto.XenditCallbackRequest true "Callback Payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /webhooks/xendit [post]
func (h *WalletHandler) HandleXenditCallback(c echo.Context) error {
	var payload dto.XenditCallbackRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid webhook payload"})
	}

	if payload.Status == "PAID" || payload.Status == "SETTLED" {
		err := h.userServ.ProcessPayment(c.Request().Context(), payload.ExternalID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to process payment"})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "callback processed"})
}

// GetBalance godoc
// @Summary Get user balance
// @Description Get the current balance of the authenticated user's wallet
// @Tags Wallet
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/wallet/balance [get]
func (h *WalletHandler) GetBalance(c echo.Context) error {
	userID := utils.GetUserID(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid token"})
	}

	balance, err := h.userServ.GetBalance(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to retrieve balance"})
	}

	return c.JSON(http.StatusOK, echo.Map{"balance": balance})
}