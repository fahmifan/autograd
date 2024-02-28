package mailer

import "context"

type Email struct {
	Subject   string
	From      string
	To        string
	Body      string
	BodyPlain string
}

type Mailer interface {
	Send(ctx context.Context, email Email) (err error)
}
