// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package resend

import (
	"fmt"
	"strings"

	resendgo "github.com/resend/resend-go/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Mailer sends emails through Resend.
type Mailer struct {
	client *resendgo.Client
	from   string
}

// NewMailer creates a Resend mailer from config file values.
func NewMailer() *Mailer {
	return &Mailer{
		client: resendgo.NewClient(strings.TrimSpace(viper.GetString("app.mailer.resend.api_key"))),
		from:   strings.TrimSpace(viper.GetString("app.mailer.from")),
	}
}

// SendInviteEmail sends the user invite email with the sign-in link.
func (m *Mailer) SendInviteEmail(to, inviteLink, platformName string) error {
	return m.sendTemplate(to, fmt.Sprintf("You're invited to %s", platformName), "invite.html", map[string]string{
		"PlatformName": platformName,
		"InviteLink":   inviteLink,
	})
}

// SendPasswordResetEmail sends the password reset email with the reset link.
func (m *Mailer) SendPasswordResetEmail(to, resetLink, platformName string) error {
	return m.sendTemplate(to, fmt.Sprintf("Reset your password - %s", platformName), "rpwd.html", map[string]string{
		"PlatformName": platformName,
		"ResetLink":    resetLink,
	})
}

// SendWelcomeEmail sends a welcome email after account registration.
func (m *Mailer) SendWelcomeEmail(to, name, loginLink, platformName string) error {
	return m.sendTemplate(to, fmt.Sprintf("Welcome to %s", platformName), "welcome.html", map[string]string{
		"PlatformName": platformName,
		"UserName":     name,
		"LoginLink":    loginLink,
	})
}

// SendVerifyEmail sends an email verification message after account registration.
func (m *Mailer) SendVerifyEmail(to, name, verifyLink, platformName string) error {
	return m.sendTemplate(to, fmt.Sprintf("Verify your email - %s", platformName), "vemail.html", map[string]string{
		"PlatformName": platformName,
		"UserName":     name,
		"VerifyLink":   verifyLink,
	})
}

// sendTemplate sends an email using a named template.
func (m *Mailer) sendTemplate(to, subject, templateName string, data any) error {
	htmlBody, err := renderTemplate(templateName, data)
	if err != nil {
		return err
	}

	_, err = m.client.Emails.Send(
		&resendgo.SendEmailRequest{
			From:    m.from,
			To:      []string{to},
			Subject: subject,
			Html:    htmlBody,
		},
	)
	if err != nil {
		return fmt.Errorf("resend send email: %w", err)
	}

	log.Info().
		Str("provider", "resend").
		Str("to", to).
		Str("subject", subject).
		Msg("Email sent")

	return nil
}
