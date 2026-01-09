package handler

import (
	"net/http"
	"rent-and-read/internal/service"
	"rent-and-read/utils"

	"github.com/labstack/echo/v4"
)

type BookHandler struct{
	serv service.BookService
}

func NewBookHandler(serv service.BookService) BookHandler{
	return BookHandler{
		serv: serv,
	}
}

// Find godoc
// @Summary List all books
// @Description Get a list of all available books
// @Tags Books
// @Produce json
// @Success 200 {array} dto.BookResponse
// @Failure 500 {object} map[string]string
// @Router /api/books [get]
func (h *BookHandler) Find(c echo.Context)error{
	books, err := h.serv.Find(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to retrieve books",
		})
	}

	return c.JSON(http.StatusOK, books)
}

// FindByUserID godoc
// @Summary List rented books
// @Description Get a list of books currently rented by the authenticated user
// @Tags Books
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.BookResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/books [get]
func (h *BookHandler) FindByUserID(c echo.Context)error{
	userID := utils.GetUserID(c)
	if userID == 0{
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid token format"})
	}
	books, err := h.serv.FindByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to retrieve books",
		})
	}

	return c.JSON(http.StatusOK, books)
}

// FindActiveByUserID godoc
// @Summary List active rented books
// @Description Get a list of books currently actively rented by the authenticated user
// @Tags Books
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.BookResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/books/active [get]
func (h *BookHandler) FindActiveByUserID(c echo.Context)error{
	userID := utils.GetUserID(c)
	if userID == 0{
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid token format"})
	}
	books, err := h.serv.FindActiveByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to retrieve active books",
		})
	}

	return c.JSON(http.StatusOK, books)
}