package main

import (
	"context"
	"log"
	"net/http"
	"rent-and-read/config"
	"rent-and-read/handler"
	"rent-and-read/internal/middleware"
	"rent-and-read/internal/repository"
	"rent-and-read/internal/service"
	"rent-and-read/pkg/xendit"

	_ "rent-and-read/docs"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Rent and Read API
// @version 1.0
// @description A robust backend for a physical book rental system.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.InitDB()
	config.InitJwt()
	config.InitEmail()

    userRepo := repository.NewPsqlUserRepository(config.DB)
    bookRepo := repository.NewPsqlBookRepository(config.DB)
    rentalRepo := repository.NewPsqlRentalRepository(config.DB)

    xClient := xendit.NewClient()
    
    var emailSvc service.EmailService
    emailSvc = service.NewSmtpEmailService(
        config.Email.Host,
        config.Email.Port,
        config.Email.Username,
        config.Email.Password,
        config.Email.From,
    )

    walletSvc := service.NewWalletService(xClient, userRepo)
    rentalSvc := service.NewRentalService(rentalRepo, bookRepo, userRepo, xClient, emailSvc)
	userSvc := service.NewUserService(userRepo)
	bookSvc := service.NewBookService(bookRepo)

    rentalHandler := handler.NewRentalHandler(rentalSvc)
    walletHandler := handler.NewWalletHandler(userSvc, walletSvc)
	bookHandler := handler.NewBookHandler(bookSvc)
	userHandler := handler.NewUserHandler(userSvc)

    e := echo.New()

    c := cron.New()
    // Run daily at 09:00 AM -> 0 9 * * *
    // Run every 5 minutes -> */5 * * * *
    _, _ = c.AddFunc("0 9 * * *", func() {
        log.Println("Running daily rental due reminders...")
        if err := rentalSvc.SendDueReminders(context.Background()); err != nil {
            log.Printf("Error sending reminders: %v", err)
        }
    })
    c.Start()

    e.GET("/", func(c echo.Context) error {
        return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
    })

    e.POST("/webhooks/xendit", walletHandler.HandleXenditCallback)
    
    e.GET("/swagger/*", echoSwagger.WrapHandler)

    e.POST("/debug/send-reminders", func(c echo.Context) error {
        if err := rentalSvc.SendDueReminders(c.Request().Context()); err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
        }
        return c.JSON(http.StatusOK, map[string]string{"message": "Reminders triggered"})
    })

    r := e.Group("/api")
	r.GET("/books", bookHandler.Find)
	r.POST("/users/register", userHandler.Register)
	r.POST("/users/login", userHandler.Login)

	protected := r.Group("")
	protected.Use(echojwt.WithConfig(middleware.JWTMiddleware()))
    
    protected.POST("/rentals", rentalHandler.CreateRental)
    protected.POST("/wallet/topup", walletHandler.TopUp)
	protected.GET("/wallet/balance", walletHandler.GetBalance)
	protected.POST("/rentals/:id/return", rentalHandler.ReturnRental)
	protected.GET("/users/books/active", bookHandler.FindActiveByUserID)
	protected.GET("/users/books", bookHandler.FindByUserID)

    e.Logger.Fatal(e.Start(":8080"))
}