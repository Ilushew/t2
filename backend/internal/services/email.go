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

// SendApplicationConfirmation отправляет подтверждение клиенту
func (s *EmailService) SendApplicationConfirmation(to, route string) error {
	subject := "Заявка принята - Udmurtia AI Route"
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<div style="background: #4F46E5; color: white; padding: 20px; text-align: center;">
			<h1>Udmurtia AI Route</h1>
		</div>
		<div style="background: #f9f9f9; padding: 30px; border: 1px solid #ddd;">
			<h2>Заявка принята!</h2>
			<p>Здравствуйте!</p>
			<p>Ваша заявка на маршрут <strong>%s</strong> успешно принята.</p>
			<p>Мы свяжемся с вами в ближайшее время для уточнения деталей.</p>
		</div>
		<div style="text-align: center; padding: 20px; color: #666; font-size: 12px;">
			<p>© 2026 Udmurtia AI Route</p>
		</div>
	</div>
</body>
</html>
	`, route)

	return s.sendEmail(to, subject, body)
}

// SendApplicationToAdmin отправляет администратору информацию о заявке
func (s *EmailService) SendApplicationToAdmin(adminEmail, routeName, applicantEmail, comment string) error {
	subject := "Заявка на маршрут: " + routeName

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<div style="background: #FF6B35; color: white; padding: 20px; text-align: center;">
			<h1>Новая заявка</h1>
		</div>
		<div style="background: #f9f9f9; padding: 30px; border: 1px solid #ddd;">
			<h2>Маршрут: %s</h2>
			
			<h3>Заявка</h3>
			<div style="background: white; padding: 20px; margin: 15px 0; white-space: pre-wrap;">%s</div>
			
			<h3>Обратная связь</h3>
			<div style="background: white; padding: 20px; margin: 15px 0;">
				<p><strong>Email:</strong> <a href="mailto:%s">%s</a></p>
			</div>
			
			<p style="color: #666; font-size: 14px;">Дата: %s</p>
		</div>
		<div style="text-align: center; padding: 20px; color: #666; font-size: 12px;">
			<p>© 2026 Udmurtia AI Route</p>
		</div>
	</div>
</body>
</html>
	`, routeName, comment, applicantEmail, applicantEmail, time.Now().Format("02.01.2006 15:04"))

	return s.sendEmail(adminEmail, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	m := mail.NewMsg()
	m.From(s.from)
	m.To(to)
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body)

	if err := s.client.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func (s *EmailService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
