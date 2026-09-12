// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package conf

import (
	"github.com/clivern/ziee/files"
	"github.com/clivern/ziee/mails"
)

var (
	// SetupIssueBody is the body of the repository setup issue.
	SetupIssueBody = files.SetupIssueBody

	// SetupPullRequestBody is the body of the starter `.ziee.yml` pull request.
	SetupPullRequestBody = files.SetupPullRequestBody

	// DefaultZieeYML is the starter policy file added on repository bootstrap.
	DefaultZieeYML = files.DefaultZieeYML

	// InviteEmail is the workspace invite email body.
	InviteEmail = mails.Invite

	// WelcomeEmail is the account welcome email body.
	WelcomeEmail = mails.Welcome

	// VerifyEmail is the email verification body.
	VerifyEmail = mails.VerifyEmail

	// ResetPwdEmail is the password reset email body.
	ResetPwdEmail = mails.ResetPwd
)
