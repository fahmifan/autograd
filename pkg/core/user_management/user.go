package user_management

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/google/uuid"
	"github.com/matcornic/hermes/v2"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type ManagedUser struct {
	ID              uuid.UUID
	Name            string
	Email           string
	Role            auth.Role
	Active          bool
	ActivationToken ActivationToken

	core.TimestampMetadata
}

func (user ManagedUser) HasToken() bool {
	return user.ActivationToken.Token != ""
}

type ActivationToken struct {
	ID        uuid.UUID
	Token     string
	ExpiresAt time.Time

	core.TimestampMetadata
}

type CreateUserRequest struct {
	NewID      uuid.UUID
	Now        time.Time
	Name       string
	Email      string
	Role       auth.Role
	NewTokenID uuid.UUID
	Token      string
}

func CreateManagedUser(req CreateUserRequest) (ManagedUser, error) {
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return ManagedUser{}, ErrInvalidEmail
	}

	if !auth.ValidRole(req.Role) {
		return ManagedUser{}, errors.New("invalid role")
	}

	if len(strings.TrimSpace(req.Name)) < 3 {
		return ManagedUser{}, errors.New("name must be at least 3 characters long")
	}

	return ManagedUser{
		ID:                req.NewID,
		Name:              req.Name,
		Email:             req.Email,
		Role:              req.Role,
		TimestampMetadata: core.NewTimestampMeta(req.Now),
		Active:            false,
		ActivationToken: ActivationToken{
			ID:                req.NewTokenID,
			Token:             req.Token,
			ExpiresAt:         req.Now.Add(tokenValidDur),
			TimestampMetadata: core.NewTimestampMeta(req.Now),
		},
	}, nil
}

type CreateAdminUserRequest struct {
	NewID      uuid.UUID
	Now        time.Time
	Name       string
	Email      string
	NewTokenID uuid.UUID
	Token      string
}

const tokenValidDur = time.Minute * 30

func CreateAdminUser(req CreateAdminUserRequest) (ManagedUser, error) {
	return CreateManagedUser(CreateUserRequest{
		NewID:      req.NewID,
		Now:        req.Now,
		Name:       req.Name,
		Email:      req.Email,
		Role:       auth.RoleAdmin,
		NewTokenID: req.NewTokenID,
		Token:      req.Token,
	})
}

type RegistrationEmail struct {
	HTMLBody      string
	PlainTextBody string
	Subject       string
	FromEmail     string
	ToEmail       string
}

type CreateRegistrationEmailRequest struct {
	SenderEmail string
	User        ManagedUser
	AppLink     string
	LogoURL     string
}

func CreateRegistrationEmail(req CreateRegistrationEmailRequest) (RegistrationEmail, error) {
	if req.SenderEmail == "" {
		return RegistrationEmail{}, errors.New("invalid sender email")
	}

	if req.User.ID == uuid.Nil {
		return RegistrationEmail{}, errors.New("invalid user id")
	}

	if req.User.Email == "" {
		return RegistrationEmail{}, errors.New("invalid user email")
	}

	if req.User.Active {
		return RegistrationEmail{}, errors.New("user already active")
	}

	hh := hermes.Hermes{
		Product: hermes.Product{
			Name: "Autograde",
			Link: req.AppLink,
			Logo: req.LogoURL,
		},
	}

	emailBody := hermes.Email{
		Body: hermes.Body{
			Name: req.User.Name,
			Intros: []string{
				"Welcome to Autograd! We're excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Autograd, please activate your account here",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Activate Account",
						Link:  createUserActivationLink(req.AppLink, req.User.ID, req.User.ActivationToken.Token),
					},
				},
			},
		},
	}

	htmlBody, err := hh.GenerateHTML(emailBody)
	if err != nil {
		return RegistrationEmail{}, fmt.Errorf("generate html body: %w", err)
	}

	txtBody, err := hh.GeneratePlainText(emailBody)
	if err != nil {
		return RegistrationEmail{}, fmt.Errorf("generate plain text body: %w", err)
	}

	regEmail := RegistrationEmail{
		Subject:       "Activate your Autograd account",
		FromEmail:     req.SenderEmail,
		ToEmail:       req.User.Email,
		HTMLBody:      htmlBody,
		PlainTextBody: txtBody,
	}

	return regEmail, nil
}

func (u ManagedUser) Activate(now time.Time, token string) (ManagedUser, error) {
	if u.Active {
		return ManagedUser{}, errors.New("user already active")
	}

	if !u.HasToken() {
		return ManagedUser{}, errors.New("token not found")
	}

	if u.ActivationToken.Token != token {
		return ManagedUser{}, errors.New("invalid token")
	}

	u.Active = true
	u.ActivationToken.ExpiresAt = now
	u.UpdatedAt = now
	return u, nil
}

func createUserActivationLink(webBaseURL string, userID uuid.UUID, activationToken string) string {
	urlVal := url.Values{}
	urlVal.Add("userID", userID.String())
	urlVal.Add("activationToken", activationToken)

	return webBaseURL + "/account-activation?" + urlVal.Encode()
}

func CheckPassword(password, passwordConfirmation string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if password != passwordConfirmation {
		return errors.New("password does not match")
	}

	return nil
}
