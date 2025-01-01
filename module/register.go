// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/resend"
	"github.com/clivern/actx0/pkg/util"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	ErrFailedCompleteRegistration = errors.New("failed complete registration")
)

// RegisterRequest is the body for signing up.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,min=4,max=255" label:"Email"`
	Name     string `json:"name" validate:"required,min=1,max=255" label:"Name"`
	Password string `json:"password" validate:"required,strong_password,min=8,max=60" label:"Password"`
}

// RegisterResult is returned on a successful registration.
type RegisterResult struct {
	User *db.User
}

// Register handles direct account registration.
type Register struct {
	UserRepository db.UserRepository
}

// NewRegister creates a Register module with the given repositories.
func NewRegister(users db.UserRepository) *Register {
	return &Register{UserRepository: users}
}

// Register creates a regular user account.
func (r *Register) Register(req *RegisterRequest) (*RegisterResult, error) {
	email := strings.TrimSpace(req.Email)
	if strings.EqualFold(email, db.BotUserEmail) {
		return nil, ErrUserEmailAlreadyExists
	}

	existing, err := r.UserRepository.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%w: check existing user: %v", ErrFailedCompleteRegistration, err)
	}
	if existing != nil {
		return nil, ErrUserEmailAlreadyExists
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: hash password: %v", ErrFailedCompleteRegistration, err)
	}

	now := time.Now().UTC()
	token := uuid.NewString()
	user := &db.User{
		Name:             strings.TrimSpace(req.Name),
		Email:            email,
		Password:         hashedPassword,
		Provider:         db.UserProviderLocal,
		Role:             db.UserRoleRegular,
		IsActive:         true,
		IsEmailVerified:  false,
		EmailVerifyToken: &token,
		LastLoginAt:      now,
		Language:         db.UserLanguageEN,
		Theme:            db.UserThemeDefault,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	err = r.UserRepository.Create(user)
	if err != nil {
		return nil, fmt.Errorf("%w: create user: %v", ErrFailedCompleteRegistration, err)
	}

	mailer := resend.NewMailer()
	err = mailer.SendVerifyEmail(
		user.Email,
		user.Name,
		fmt.Sprintf("%s/verify-email/%s", viper.GetString("app.url"), token),
		viper.GetString("app.name"),
	)
	if err != nil {
		log.Warn().
			Err(err).
			Str("userId", user.Id.String()).
			Str("email", user.Email).
			Msg("Failed to send verify email; registration completed")
	}

	err = mailer.SendWelcomeEmail(
		user.Email,
		user.Name,
		fmt.Sprintf("%s/login", viper.GetString("app.url")),
		viper.GetString("app.name"),
	)
	if err != nil {
		log.Warn().
			Err(err).
			Str("userId", user.Id.String()).
			Str("email", user.Email).
			Msg("Failed to send welcome email; registration completed")
	}

	return &RegisterResult{User: user}, nil
}
