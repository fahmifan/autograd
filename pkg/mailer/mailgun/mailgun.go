package mailgun

import (
	"context"
	"time"

	"github.com/fahmifan/autograd/pkg/mailer"
	"github.com/mailgun/mailgun-go/v4"
)

type MailgunClient interface {
	Send(ctx context.Context, m *mailgun.Message) (string, string, error)
	NewMessage(from, subject, text string, to ...string) *mailgun.Message
}

type MailgunTransporter struct {
	client MailgunClient
}

func (m *MailgunTransporter) Send(ctx context.Context, email mailer.Email) (err error) {
	message := m.client.NewMessage(email.From, email.Subject, email.BodyPlain, email.To)
	message.SetHtml(email.Body)
	ctx2, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, err = m.client.Send(ctx2, message)
	return err
}
