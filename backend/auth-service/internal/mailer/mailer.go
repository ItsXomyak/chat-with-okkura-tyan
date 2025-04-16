package mailer

import (
	"fmt"
	"net/smtp"

	"authforge/config"
	"authforge/internal/logger"
)

type Mailer interface {
	SendConfirmationEmail(to, token string) error
	SendPasswordResetEmail(to, token string) error
}

type smtpMailer struct {
	cfg *config.Config
}

func NewSMTPMailer(cfg *config.Config) Mailer {
	return &smtpMailer{cfg: cfg}
}

func (m *smtpMailer) SendConfirmationEmail(to, token string) error {
	subject := "Account Confirmation"
	confirmationURL := fmt.Sprintf("http://localhost:8080/api/v1/auth/confirm?token=%s", token)
	body := fmt.Sprintf("Please confirm your account by clicking the link:\n%s", confirmationURL)
	logger.Info("Sending confirmation email to ", to)
	err := m.sendMail(to, subject, body)
	if err != nil {
		logger.Error("Error sending confirmation email to ", to, ": ", err)
	} else {
		logger.Info("Confirmation email sent to ", to)
	}
	return err
}

func (m *smtpMailer) SendPasswordResetEmail(to, token string) error {
	subject := "Password Reset"
	resetInstructions := fmt.Sprintf(
		"To reset your password, send a POST request to endpoint http://localhost:8080/api/v1/auth/password-reset-confirm with the following JSON body:\n{\n\t\"token\": \"%s\",\n\t\"newPassword\": \"<your new password>\"\n}",
		token,
	)
	logger.Info("Sending password reset email to ", to)
	err := m.sendMail(to, subject, resetInstructions)
	if err != nil {
		logger.Error("Error sending password reset email to ", to, ": ", err)
	} else {
		logger.Info("Password reset email sent to ", to)
	}
	return err
}

func (m *smtpMailer) sendMail(to, subject, body string) error {
	from := m.cfg.SMTPUsername
	password := m.cfg.SMTPPassword
	host := m.cfg.SMTPHost
	port := m.cfg.SMTPPort

	addr := fmt.Sprintf("%s:%d", host, port)
	logger.Info("Connecting to SMTP server at ", addr)

	auth := smtp.PlainAuth("", from, password, host)

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			body + "\r\n",
	)

	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		logger.Error("SMTP SendMail error: ", err)
	}
	return err
}
