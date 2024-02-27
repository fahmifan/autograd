package smtp

import (
	"context"
	_ "embed"

	"github.com/fahmifan/autograd/pkg/mailer"
	"gopkg.in/gomail.v2"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

type SMTP struct {
	mail *gomail.Dialer
	cfg  *Config
}

func NewSmtpClient(cfg *Config) (smtp *SMTP, err error) {
	smtp = &SMTP{
		cfg: cfg,
		mail: gomail.NewDialer(
			cfg.Host,
			cfg.Port,
			cfg.Username,
			cfg.Password,
		),
	}

	closer, err := smtp.mail.Dial()
	if err != nil {
		return nil, err
	}
	closer.Close()

	return smtp, nil
}

func (m *SMTP) Send(ctx context.Context, email mailer.Email) (err error) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", email.From)
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)
	msg.SetBody("text/html", email.Body)
	msg.AddAlternative("text/plain", email.BodyPlain)

	return m.mail.DialAndSend(msg)
}
