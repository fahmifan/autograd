package user_management_cmd

import (
	"context"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/user_management"
	"github.com/fahmifan/autograd/pkg/jobqueue"
	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/fahmifan/autograd/pkg/mailer"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const JobSendEmail jobqueue.JobType = "send_email"

type SendRegistrationEmailHandler struct {
	*core.Ctx
}

type SendRegistrationEmailPayload struct {
	UserID uuid.UUID
}

func (handler *SendRegistrationEmailHandler) JobType() jobqueue.JobType {
	return JobSendEmail
}

func (handler *SendRegistrationEmailHandler) Handle(ctx context.Context, tx *gorm.DB, payload jobqueue.Payload) error {
	req := SendRegistrationEmailPayload{}
	err := jobqueue.UnmarshalPayload(payload, &req)
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "SendRegistrationEmailHandler: Handle: json.Unmarshal")
	}

	user, err := user_management.ManagedUserReader{}.FindUserByID(ctx, tx, req.UserID)
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "SendRegistrationEmailHandler: Handle: FindUserByID")
	}

	regEmail, err := user_management.CreateRegistrationEmail(user_management.CreateRegistrationEmailRequest{
		User:        user,
		SenderEmail: handler.SenderEmail,
		AppLink:     handler.AppLink,
		LogoURL:     handler.LogoURL,
	})
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "SendRegistrationEmailHandler: Handle: CreateRegistrationEmail")
	}

	err = handler.Ctx.Mailer.Send(ctx, mailer.Email{
		Subject:   regEmail.Subject,
		From:      regEmail.FromEmail,
		To:        regEmail.ToEmail,
		Body:      regEmail.HTMLBody,
		BodyPlain: regEmail.PlainTextBody,
	})
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "SendRegistrationEmailHandler: Handle: Mailer.Send")
	}

	return nil
}
