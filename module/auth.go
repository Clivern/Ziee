// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/actx0/ziee/conf"
	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/pkg/util"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrAccountNotActive       = errors.New("account not active")
	ErrMaintenanceModeEnabled = errors.New("maintenance mode enabled")
	ErrFailedCheckMaintenance = errors.New("failed check maintenance")
	ErrFailedCreateSession    = errors.New("failed create session")

	ErrUserNotFound           = errors.New("user not found")
	ErrUserEmailAlreadyExists = errors.New("user with this email already exists")
	ErrResetTokenInvalid      = errors.New("reset link invalid or expired")
	ErrFailedForgotPassword   = errors.New("failed process forgot password")
	ErrFailedResetPassword    = errors.New("failed reset password")

	ErrFailedRevokeSessions = errors.New("failed revoke sessions")
)

// LoginRequest is the body you send to log in.
type LoginRequest struct {
	Email      string `json:"email" validate:"required,email" label:"Email"`
	Password   string `json:"password" validate:"required" label:"Password"`
	RememberMe bool   `json:"rememberMe" validate:"omitempty,boolean" label:"Remember Me"`
}

// LoginResult is returned on a successful login.
type LoginResult struct {
	User          *db.User
	Session       *db.Session
	CookieOptions *util.CookieOptions
}

// ForgotPasswordRequest is the body for requesting a password reset (just email).
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,min=4,max=255" label:"Email"`
}

// ResetPasswordRequest is the body for resetting a password (token + new password).
type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required,min=1" label:"Token"`
	Password string `json:"password" validate:"required,strong_password,min=8,max=60" label:"Password"`
}

// Mailer sends the transactional emails used by the auth flows.
type Mailer interface {
	SendPasswordResetEmail(to, resetLink, pname string) error
}

// Auth handles user authentication, sessions, and password resets.
type Auth struct {
	UserRepository               db.UserRepository
	SessionRepository            db.SessionRepository
	ConfigRepository             db.ConfigRepository
	PasswordResetTokenRepository db.PasswordResetTokenRepository
	Mailer                       Mailer
}

// NewAuth creates an Auth module with the given repositories.
func NewAuth(
	users db.UserRepository,
	sessions db.SessionRepository,
	config db.ConfigRepository,
	resetTokens db.PasswordResetTokenRepository,
	mailer Mailer,
) *Auth {
	return &Auth{
		UserRepository:               users,
		SessionRepository:            sessions,
		ConfigRepository:             config,
		PasswordResetTokenRepository: resetTokens,
		Mailer:                       mailer,
	}
}

// Login authenticates the user and creates a session.
func (a *Auth) Login(ctx context.Context, req *LoginRequest) (*LoginResult, error) {
	user, err := a.UserRepository.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	if user == nil || !util.ComparePassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}
	if db.IsBotUser(user.Id) {
		return nil, ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, ErrAccountNotActive
	}

	maintenance, err := a.ConfigRepository.Get("maintenance_mode")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCheckMaintenance, err)
	}
	if maintenance != nil && maintenance.Value == "1" && user.Role != db.UserRoleAdmin {
		return nil, ErrMaintenanceModeEnabled
	}

	err = a.UserRepository.UpdateLastLogin(user.Id)
	if err != nil {
		return nil, fmt.Errorf("update last login: %w", err)
	}

	duration := lo.Ternary(req.RememberMe, conf.RememberMeSessionDuration, conf.DefaultSessionDuration)

	token, err := util.GenerateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateSession, err)
	}

	session := &db.Session{
		Token:     token,
		UserId:    user.Id,
		ExpiresAt: time.Now().UTC().Add(duration),
	}

	if appCtx, ok := AppCtxFrom(ctx); ok {
		if lo.IsNotEmpty(appCtx.IPAddress) {
			session.IPAddress = &appCtx.IPAddress
		}
		if lo.IsNotEmpty(appCtx.UserAgent) {
			session.UserAgent = &appCtx.UserAgent
		}
	}

	err = a.SessionRepository.Create(session)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateSession, err)
	}

	log.Info().
		Str("userId", user.Id.String()).
		Str("sessionId", session.Id.String()).
		Msg("User logged in")

	UserLoggedIn.Emit(ctx, UserLoggedInData{
		UserID:    user.Id,
		SessionID: session.Id,
	})

	cookieOptions := lo.Ternary(
		strings.HasPrefix(viper.GetString("app.url"), "https://"),
		util.SecureCookieOptions(),
		util.DefaultCookieOptions(),
	)
	cookieOptions.MaxAge = lo.Ternary(req.RememberMe, int(conf.RememberMeSessionDuration/time.Second), 0)

	return &LoginResult{
		User:          user,
		Session:       session,
		CookieOptions: cookieOptions,
	}, nil
}

// Logout revokes all sessions for a user.
func (a *Auth) Logout(ctx context.Context, userId db.Id) error {
	count, err := a.SessionRepository.DeleteExpired()
	if err == nil && count > 0 {
		log.Info().
			Int64("count", count).
			Msg("Expired sessions deleted")
	}

	err = a.SessionRepository.DeleteByUserId(userId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedRevokeSessions, err)
	}

	log.Info().
		Str("userId", userId.String()).
		Msg("User logged out")

	UserLoggedOut.Emit(ctx, UserLoggedOutData{UserID: userId})

	return nil
}

// ForgotPassword creates a password reset token and emails the reset link.
func (a *Auth) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	user, err := a.UserRepository.GetByEmail(req.Email)
	if err != nil {
		log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("Failed to look up user for forgot password")
		return nil
	}
	if user == nil {
		return nil
	}
	if db.IsBotUser(user.Id) {
		return nil
	}

	token, _ := uuid.NewRandom()
	resetToken := &db.PasswordResetToken{
		UserId:    user.Id,
		Token:     token.String(),
		ExpiresAt: time.Now().UTC().Add(conf.PasswordResetTokenTTL),
	}

	err = a.PasswordResetTokenRepository.Create(resetToken)
	if err != nil {
		return fmt.Errorf("%w: create reset token: %v", ErrFailedForgotPassword, err)
	}

	purl := strings.TrimRight(viper.GetString("app.url"), "/")
	pname := viper.GetString("app.name")

	err = a.Mailer.SendPasswordResetEmail(
		user.Email,
		fmt.Sprintf("%s/reset-password/%s", purl, url.PathEscape(token.String())),
		pname,
	)
	if err != nil {
		log.Warn().
			Err(err).
			Str("email", user.Email).
			Msg("Failed to send password reset email; token was created")
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Password reset token created")

	UserForgotPassword.Emit(ctx, UserForgotPasswordData{UserID: user.Id})

	return nil
}

// ResetPassword sets a new password using a valid reset token.
func (a *Auth) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	tok, err := a.PasswordResetTokenRepository.GetByToken(req.Token)
	if err != nil {
		log.Error().
			Err(err).
			Str("token", req.Token).
			Msg("Failed to look up password reset token")
		return ErrResetTokenInvalid
	}
	if tok == nil || time.Now().UTC().After(tok.ExpiresAt) {
		return ErrResetTokenInvalid
	}

	user, err := a.UserRepository.GetById(tok.UserId)
	if err != nil || user == nil {
		return ErrResetTokenInvalid
	}

	hpass, err := util.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("%w: hash password: %v", ErrFailedResetPassword, err)
	}

	user.Password = hpass
	err = a.UserRepository.Update(user)
	if err != nil {
		return fmt.Errorf("%w: update user: %v", ErrFailedResetPassword, err)
	}

	err = a.PasswordResetTokenRepository.DeleteByToken(req.Token)
	if err != nil {
		log.Error().
			Err(err).
			Str("token", req.Token).
			Str("userId", tok.UserId.String()).
			Msg("Failed to delete password reset token")
	}

	log.Info().
		Str("userId", tok.UserId.String()).
		Msg("Password reset completed")

	UserResetPassword.Emit(ctx, UserResetPasswordData{UserID: tok.UserId})

	return nil
}
