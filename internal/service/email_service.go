package service

import (
	"fmt"
	"net/smtp"
	"time"
)

type ConsoleEmailService struct{}

func NewConsoleEmailService() *ConsoleEmailService {
	return &ConsoleEmailService{}
}

type SmtpEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSmtpEmailService(host, port, username, password, from string) *SmtpEmailService {
	return &SmtpEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SmtpEmailService) SendRentalDueReminder(toEmail, name, bookTitle string, dueDate time.Time) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = toEmail
	headers["Subject"] = "Rental Due Reminder"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + fmt.Sprintf("<html><body><p>Hi %s,</p><p>This is a reminder that your book <b>%s</b> is due on %s.</p></body></html>", name, bookTitle, dueDate.Format("2006-01-02"))

	return smtp.SendMail(addr, auth, s.from, []string{toEmail}, []byte(message))
}