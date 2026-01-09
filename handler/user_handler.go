package handler

import (
	"errors"
	"net/http"
	"rent-and-read/internal/dto"
	"rent-and-read/internal/models"
	"rent-and-read/internal/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct{
	serv service.UserService
}

func NewUserHandler(serv service.UserService) UserHandler{
	return UserHandler{
		serv: serv,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Router /api/users/register [post]
func (h *UserHandler) Register (c echo.Context) error{
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil{
		return c.JSON(http.StatusBadRequest, echo.Map{"message": models.ErrBadUserRequest})
	}

	err := h.serv.Register(c.Request().Context(), req)
	if err != nil{
		if errors.Is(err, models.ErrUserNotFound){
			return c.JSON(http.StatusNotFound, echo.Map{"message": models.ErrUserNotFound})
		}else if errors.Is(err, models.ErrUserExist){
			return c.JSON(http.StatusNotFound, echo.Map{"message": models.ErrUserExist})
		}else if errors.Is(err, models.ErrBadUserRequest){
			return c.JSON(http.StatusNotFound, echo.Map{"message": models.ErrBadUserRequest})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error", 
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "succesfully created a user",
		"user": dto.UserResponse{
			Name: req.Name,
			Email: req.Email,
		},
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/login [post]
func (h *UserHandler) Login (c echo.Context) error{
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil{
		return c.JSON(http.StatusBadRequest, echo.Map{"message": models.ErrBadUserRequest})
	}

	token, err := h.serv.Login(c.Request().Context(), req)
	if err != nil{
		if errors.Is(err, models.ErrUnauthorized){
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": models.ErrUnauthorized,
			})
		} else if errors.Is(err, models.ErrBadUserRequest) {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": models.ErrBadUserRequest})
		}
		
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error", 
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}