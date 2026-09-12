// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package conf

import (
	"time"

	"github.com/clivern/ziee/files"
	"github.com/clivern/ziee/mails"
)

const (
	// SlowRequestThreshold is the latency at or above which a request is counted as slow.
	SlowRequestThreshold = time.Second

	// DefaultSessionDuration is the auth session TTL without remember-me.
	DefaultSessionDuration = 7 * 24 * time.Hour

	// RememberMeSessionDuration is the auth session TTL with remember-me.
	RememberMeSessionDuration = 30 * 24 * time.Hour

	// InviteExpiry is how long a workspace invite stays valid.
	InviteExpiry = 7 * 24 * time.Hour

	// MaxUploadBytes is the maximum multipart upload size.
	MaxUploadBytes = 2 * 1024 * 1024

	// DefaultSearchLimit is used when a knowledge search request omits limit.
	DefaultSearchLimit = 10

	// ConfigSyncCheckName is the GitHub check run created when syncing `.ziee.yml`.
	ConfigSyncCheckName = "Ziee configurations sync"

	// SetupIssueTitle is the GitHub issue opened when Ziee is installed on a repository.
	SetupIssueTitle = "Set up Ziee for this repository"

	// SetupPullRequestTitle is the pull request that adds a starter `.ziee.yml`.
	SetupPullRequestTitle = "Add .ziee.yml"
)

// SetupIssueBody is the body of the repository setup issue.
var SetupIssueBody = files.SetupIssueBody

// SetupPullRequestBody is the body of the starter `.ziee.yml` pull request.
var SetupPullRequestBody = files.SetupPullRequestBody

// DefaultZieeYML is the starter policy file added on repository bootstrap.
var DefaultZieeYML = files.DefaultZieeYML

// InviteEmail is the workspace invite email body.
var InviteEmail = mails.Invite

// WelcomeEmail is the account welcome email body.
var WelcomeEmail = mails.Welcome

// VerifyEmail is the email verification body.
var VerifyEmail = mails.VerifyEmail

// ResetPwdEmail is the password reset email body.
var ResetPwdEmail = mails.ResetPwd
