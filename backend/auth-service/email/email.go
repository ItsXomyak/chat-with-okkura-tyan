package email

import (
	"auth-service/config"
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

func (es *EmailService) SendVerificationEmail(to, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.cfg.SMTPUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Подтверждение регистрации")
	
	verificationLink := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)
	body := fmt.Sprintf("Для подтверждения регистрации перейдите по ссылке: %s", verificationLink)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(es.cfg.SMTPHost, es.cfg.SMTPPort, es.cfg.SMTPUser, es.cfg.SMTPPass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
} 