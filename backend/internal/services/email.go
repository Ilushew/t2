package services

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/template"
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
	if err != nil {
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
	subject := "Заявка принята - Маршруты по Удмуртии"

	const tmplStr = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Заявка принята</title>
    <style>
        @media only screen and (max-width: 600px) {
            .container { width: 100% !important; }
            .content { padding: 28px 16px !important; }
        }
    </style>
</head>
<body style="margin: 0; padding: 0; background-color: #0a0a0a; font-family: Inter, -apple-system, BlinkMacSystemFont, Arial, sans-serif; -webkit-font-smoothing: antialiased;">
    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="background-color: #0a0a0a;">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" class="container" style="max-width: 600px; width: 100%; background-color: #141414; border-radius: 24px; overflow: hidden; border: 1px solid #2a2a2a; box-shadow: 0 20px 60px rgba(0,0,0,0.6);">
                    <!-- Header -->
                    <tr>
                        <td align="center" style="background-color: #ff40d0; background: linear-gradient(135deg, #ff40d0 0%, #d600a0 100%); padding: 32px 20px;">
                            <h1 style="margin: 0; font-size: 20px; font-weight: 800; color: #ffffff; text-transform: uppercase; letter-spacing: 2.5px;">Маршруты по Удмуртии</h1>
                        </td>
                    </tr>
                    
                    <!-- Content -->
                    <tr>
                        <td class="content" style="padding: 36px 32px; background-color: #141414;">
                            <h2 style="margin: 0 0 24px; font-size: 20px; font-weight: 700; color: #ffffff; text-transform: uppercase; letter-spacing: 1px;">Заявка принята!</h2>
                            
                            <p style="margin: 0 0 16px; color: #cccccc; font-size: 15px; line-height: 1.6;">Здравствуйте!</p>
                            
                            <p style="margin: 0 0 16px; color: #cccccc; font-size: 15px; line-height: 1.6;">
                                Ваша заявка на маршрут успешно принята.
                            </p>
                            
                            <p style="margin: 0 0 32px; color: #cccccc; font-size: 15px; line-height: 1.6;">
                                Мы свяжемся с вами в ближайшее время для уточнения деталей. Спасибо, что выбираете Удмуртию!
                            </p>

                            <div style="border-top: 1px solid #2a2a2a; margin: 24px 0;"></div>
                            <p style="margin: 0; color: #888888; font-size: 13px; line-height: 1.5;">
                                С уважением,<br>
                                <span style="color: #ff40d0; font-weight: 600;">Команда Udmurtia AI Route</span>
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td align="center" style="padding: 24px 20px; background-color: #0f0f0f; border-top: 1px solid #222222;">
                            <p style="margin: 0; color: #555555; font-size: 11px; font-family: Inter, Arial, sans-serif; letter-spacing: 0.5px;">© 2026 Udmurtia AI Route. Хакатон проект.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

	// Парсинг шаблона
	tmpl, err := template.New("confirmation_email").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parse confirmation template: %w", err)
	}

	// Данные для подстановки
	data := struct {
		Route string
	}{
		Route: route,
	}

	// Рендер в буфер
	var bodyBuf bytes.Buffer
	if err := tmpl.Execute(&bodyBuf, data); err != nil {
		return fmt.Errorf("render confirmation template: %w", err)
	}

	return s.sendEmail(to, subject, bodyBuf.String())
}

// SendApplicationToAdmin отправляет администратору информацию о заявке
func (s *EmailService) SendApplicationToAdmin(adminEmail, routeName, applicantEmail, comment string) error {
	subject := "Заявка на маршрут"

	// Шаблон письма (html/template автоматически экранирует HTML/XSS)
	const tmplStr = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Новая заявка</title>
    <style>
        @media only screen and (max-width: 600px) {
            .container { width: 100% !important; }
            .content { padding: 24px 16px !important; }
            .detail-box { padding: 16px !important; }
        }
    </style>
</head>
<body style="margin: 0; padding: 0; background-color: #0a0a0a; font-family: Inter, -apple-system, BlinkMacSystemFont, Arial, sans-serif; -webkit-font-smoothing: antialiased;">
    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="background-color: #0a0a0a;">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" class="container" style="max-width: 600px; width: 100%; background-color: #141414; border-radius: 24px; overflow: hidden; border: 1px solid #2a2a2a; box-shadow: 0 20px 60px rgba(0,0,0,0.6);">
                    <!-- Header -->
                    <tr>
                        <td align="center" style="background-color: #ff40d0; background: linear-gradient(135deg, #ff40d0 0%, #d600a0 100%); padding: 32px 20px;">
                            <h1 style="margin: 0; font-size: 20px; font-weight: 800; color: #ffffff; text-transform: uppercase; letter-spacing: 2.5px;"> Заявка </h1>
                        </td>
                    </tr>
                    
                    <!-- Content -->
                    <tr>
                        <td class="content" style="padding: 36px 32px; background-color: #141414;">
                            <p style="margin: 0 0 10px; font-size: 13px; font-weight: 600; color: #aaaaaa; text-transform: uppercase; letter-spacing: 0.5px;">Текст заявки:</p>
                            <table role="presentation" width="100%" style="margin-bottom: 24px;">
                                <tr>
                                    <td class="detail-box" style="padding: 20px; background-color: #1e1e1e; border: 1px solid #333333; border-radius: 14px; color: #cccccc; font-size: 15px; line-height: 1.6; white-space: pre-wrap; word-wrap: break-word;">{{.Comment}}</td>
                                </tr>
                            </table>

                            <p style="margin: 0 0 10px; font-size: 13px; font-weight: 600; color: #aaaaaa; text-transform: uppercase; letter-spacing: 0.5px;">Контакт:</p>
                            <table role="presentation" width="100%" style="margin-bottom: 24px;">
                                <tr>
                                    <td class="detail-box" style="padding: 16px 20px; background-color: #1e1e1e; border: 1px solid #333333; border-radius: 12px;">
                                        <span style="color: #888888; font-size: 12px; display: block; margin-bottom: 4px; text-transform: uppercase;">Email для связи:</span>
                                        <a href="mailto:{{.Email}}" style="color: #ff40d0; font-size: 15px; font-weight: 600; text-decoration: none;">{{.Email}}</a>
                                    </td>
                                </tr>
                            </table>

                            <div style="border-top: 1px solid #2a2a2a; margin: 20px 0;"></div>
                            <p style="margin: 0; color: #666666; font-size: 13px; text-align: right;">Дата: {{.Date}}</p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td align="center" style="padding: 24px 20px; background-color: #0f0f0f; border-top: 1px solid #222222;">
                            <p style="margin: 0; color: #555555; font-size: 11px; letter-spacing: 0.5px;">© 2026 Udmurtia AI Route. Хакатон проект.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

	// Парсинг шаблона
	tmpl, err := template.New("application_email").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parse email template: %w", err)
	}

	// Данные для подстановки
	data := struct {
		RouteName string
		Comment   string
		Email     string
		Date      string
	}{
		RouteName: routeName,
		Comment:   comment,
		Email:     applicantEmail,
		Date:      time.Now().Format("02.01.2006 15:04"),
	}

	// Рендер в буфер
	var bodyBuf bytes.Buffer
	if err := tmpl.Execute(&bodyBuf, data); err != nil {
		return fmt.Errorf("render email template: %w", err)
	}

	return s.sendEmail(adminEmail, subject, bodyBuf.String())
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
