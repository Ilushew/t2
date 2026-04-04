package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/wneessen/go-mail"
)

type EmailService struct {
	client *mail.Client
	from   string
}

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewEmailService(cfg EmailConfig) (*EmailService, error) {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil{
		return nil, fmt.Errorf("")
	}
	client, err := mail.NewClient(
		cfg.Host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.Username),
		mail.WithPassword(cfg.Password),
		mail.WithSSL(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create email client: %w", err)
	}
	return &EmailService{
		client: client,
		from:   cfg.From,
	}, nil
}

func (s *EmailService) SendVerificationCode(to, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	m := mail.NewMsg()
	m.From(s.from)
	m.To(to)
	m.Subject("Код подтверждения - Udmurtia AI Route")

	m.SetBodyString(mail.TypeTextHTML, getVerificationEmailHTML(code))

	if err := s.client.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("fail to send email: %w", err)
	}
	return nil
}

func (s *EmailService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
