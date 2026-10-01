// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/github/oauth"
	"github.com/clivern/ziee/pkg/resend"
	"github.com/clivern/ziee/pkg/storage"

	"github.com/spf13/viper"
)

// Instance is the process-wide API wired at server startup.
var Instance *API

// API holds shared application modules for HTTP handlers.
type API struct {
	Auth         *module.Auth
	Profile      *module.Profile
	Setup        *module.Setup
	Settings     *module.Settings
	APIKey       *module.APIKey
	Me           *module.Me
	Workspace    *module.Workspace
	Invite       *module.Invite
	Access       *module.Access
	Billing      *module.Billing
	Stats        *module.Stats
	Audit        *module.Audit
	Document     *module.Document
	Installation *module.Installation
	Repository   *module.Repository
	OAuth        *oauth.OAuth

	Documents      db.DocumentRepository
	WorkspaceUsers db.WorkspaceUserRepository
	Usage          db.UsageRepository
	Subscriptions  db.SubscriptionRepository
	Store          storage.Store
}

// New wires repositories and modules once for the process.
func New() *API {
	conn := db.GetDB()

	store, err := storage.New()
	if err != nil {
		panic(err)
	}

	users := db.NewUserRepository(conn)
	sessions := db.NewSessionRepository(conn)
	configs := db.NewConfigRepository(conn)
	workspaces := db.NewWorkspaceRepository(conn)
	workspaceUsers := db.NewWorkspaceUserRepository(conn)
	invites := db.NewUserInviteRepository(conn)
	apiKeys := db.NewAPIKeyRepository(conn)
	accessKeys := db.NewAccessKeyRepository(conn)
	subscriptions := db.NewSubscriptionRepository(conn)
	purchases := db.NewTokenPurchaseRepository(conn)
	stats := db.NewWorkspaceStatsRepository(conn)
	audits := db.NewAuditEventRepository(conn)
	documents := db.NewDocumentRepository(conn)
	installations := db.NewGitHubInstallationRepository(conn)
	repos := db.NewRepositoriesRepository(conn)
	repoMeta := db.NewRepositoryMetaRepository(conn)
	spam := db.NewRepositorySpamUserRepository(conn)
	usage := db.NewUsageRepository(conn)

	a := &API{
		Auth:         module.NewAuth(users, sessions, configs),
		Profile:      module.NewProfile(users),
		Setup:        module.NewSetup(configs, users),
		Settings:     module.NewSettings(configs),
		APIKey:       module.NewAPIKey(apiKeys),
		Me:           module.NewMe(apiKeys, users, accessKeys),
		Workspace:    module.NewWorkspace(workspaces, workspaceUsers, subscriptions, users),
		Invite:       module.NewInvite(invites, users, configs, workspaces, workspaceUsers, resend.NewMailer()),
		Access:       module.NewAccess(accessKeys, workspaces),
		Billing:      module.NewBilling(workspaces, subscriptions, purchases, module.Usage{}),
		Stats:        module.NewStats(workspaces, stats),
		Audit:        module.NewAudit(audits, workspaces),
		Document:     module.NewDocument(documents, workspaces, store),
		Installation: module.NewInstallation(installations, repos),
		Repository:   module.NewRepository(repos, repoMeta, spam),
		OAuth: oauth.NewOAuth(oauth.OAuthConfig{
			ClientID:     viper.GetString("app.oauth.github.client_id"),
			ClientSecret: viper.GetString("app.oauth.github.client_secret"),
			RedirectURL:  viper.GetString("app.oauth.github.redirect_url"),
			Scopes:       []string{"read:user", "user:email"},
			AllowSignup:  true,
		}),
		Documents:      documents,
		WorkspaceUsers: workspaceUsers,
		Usage:          usage,
		Subscriptions:  subscriptions,
		Store:          store,
	}

	Instance = a

	return a
}
