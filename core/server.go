// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/clivern/ziee/api"
	"github.com/clivern/ziee/conf"
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/middleware"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/github"

	"github.com/go-chi/chi/v5"
	cmid "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Setup creates and configures the HTTP server
func SetupServer(Static embed.FS) http.Handler {
	r := chi.NewRouter()

	r.Use(cmid.Recoverer)
	r.Use(cmid.Timeout(
		time.Duration(viper.GetInt("app.timeout")) * time.Second,
	))
	r.Use(middleware.RequestId)
	r.Use(middleware.AppContext)
	r.Use(middleware.PrometheusMiddleware)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestSizeLimit(
		int64(10 * 1024 * 1024),
	))
	r.Use(middleware.Auth())
	r.Get("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) { // suppress browser favicon requests
		w.WriteHeader(http.StatusNoContent)
	})

	r.Group(func(r chi.Router) { // public — no auth required
		r.Get("/api/v1/public/_health", api.HealthAction)                                   // liveness probe
		r.Get("/api/v1/public/_ready", api.ReadyAction)                                     // readiness probe (db, deps)
		r.Post("/api/v1/public/action/setup", api.SetupAction)                              // initial app setup
		r.Get("/api/v1/public/action/setup/status", api.SetupStatusAction)                  // setup completion status
		r.Post("/api/v1/public/action/logout", api.LogoutAction)                            // user logout
		r.Get("/api/v1/public/action/oauth/github", api.GitHubOAuthStartAction)             // start GitHub OAuth
		r.Get("/api/v1/public/action/oauth/github/callback", api.GitHubOAuthCallbackAction) // GitHub OAuth callback
		r.Post("/api/v1/public/action/github/webhook", api.GitHubWebhookAction)             // GitHub App webhook
		if conf.IsSaaS() {
			r.Post("/api/v1/public/action/stripe/webhook", api.StripeWebhookAction) // Stripe billing webhook
		}
	})
	r.Get("/api/v1/me", api.GetMeAction) // current authenticated user
	r.Group(func(r chi.Router) {         // user profile and workspace invites
		r.Use(middleware.Protect(middleware.Config{Roles: []string{db.UserRoleAdmin, db.UserRoleRegular}}))
		r.Get("/api/v1/action/profile", api.GetProfileAction)                                                     // get user profile
		r.Put("/api/v1/action/profile", api.UpdateProfileAction)                                                  // update user profile
		r.Get("/api/v1/action/invites", api.ListUserInvitesAction)                                                // list pending workspace invites
		r.Get("/api/v1/action/github/installations", api.ListGitHubInstallationsAction)                           // list pending GitHub App installations for the user
		r.Post("/api/v1/action/github/installations/{installationId}/attach", api.AttachGitHubInstallationAction) // attach a GitHub App installation to a workspace
		r.Get("/api/v1/action/invite-by-token/{token}", api.GetAuthenticatedUserInviteByTokenAction)              // get invite details by token
		r.Post("/api/v1/action/accept-invite/{token}", api.AcceptUserInviteByTokenAction)                         // accept workspace invite
		r.Post("/api/v1/action/reject-invite/{token}", api.RejectUserInviteByTokenAction)                         // reject workspace invite
	})
	r.Group(func(r chi.Router) { // app settings — admin only
		r.Use(middleware.Protect(middleware.Config{Roles: []string{db.UserRoleAdmin}}))
		r.Put("/api/v1/action/settings", api.UpdateSettingsAction) // update app settings
		r.Get("/api/v1/action/settings", api.GetSettingsAction)    // get app settings
	})
	r.Group(func(r chi.Router) { // user API keys
		r.Use(middleware.Protect(middleware.Config{Roles: []string{db.UserRoleAdmin, db.UserRoleRegular}}))
		r.Post("/api/v1/apiKeys", api.CreateUserAPIKeyAction)              // create user API key
		r.Get("/api/v1/apiKeys", api.ListUserAPIKeysAction)                // list user API keys
		r.Get("/api/v1/apiKeys/{apiKeyId}", api.GetUserAPIKeyAction)       // get user API key
		r.Delete("/api/v1/apiKeys/{apiKeyId}", api.DeleteUserAPIKeyAction) // delete user API key
	})
	r.Group(func(r chi.Router) { // workspaces
		r.Use(middleware.Protect(middleware.Config{Roles: []string{db.UserRoleAdmin, db.UserRoleRegular}}))
		r.Post("/api/v1/workspaces", api.CreateWorkspaceAction) // create workspace
		r.Get("/api/v1/workspaces", api.ListWorkspacesAction)   // list user workspaces
	})

	r.Route("/api/v1/workspaces/{workspaceId}", func(r chi.Router) { // workspace-scoped resources
		r.Use(middleware.Protect(middleware.Config{Workspace: true}))

		r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanGetWorkspace})).Get("/", api.GetWorkspaceAction)          // get workspace
		r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanUpdateWorkspace})).Put("/", api.UpdateWorkspaceAction)    // update workspace
		r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanDeleteWorkspace})).Delete("/", api.DeleteWorkspaceAction) // delete workspace

		r.Route("/invites", func(r chi.Router) {
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanInviteMember})).Post("/", api.CreateInviteAction)             // create member invite
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanListWorkspaceInvites})).Get("/", api.ListInvitesAction)       // list member invites
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanGetWorkspaceInvite})).Get("/{inviteId}", api.GetInviteAction) // get member invite
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanRemoveInvite})).Delete("/{inviteId}", api.DeleteInviteAction) // revoke member invite
		})

		r.Route("/members", func(r chi.Router) {
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanListWorkspaceMembers})).Get("/", api.ListWorkspaceMembersAction)                // list workspace members
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanUpdateMemberRole})).Put("/{memberUserId}", api.UpdateWorkspaceMemberRoleAction) // update member role
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanRemoveMember})).Delete("/{memberUserId}", api.DeleteWorkspaceMemberAction)      // remove member
		})

		r.Route("/keys", func(r chi.Router) {
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanCreateWorkspaceAccessKey})).Post("/", api.CreateWorkspaceAccessKeyAction)          // create workspace access key
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanListWorkspaceAccessKeys})).Get("/", api.ListWorkspaceAccessKeysAction)             // list workspace access keys
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanGetWorkspaceAccessKey})).Get("/{keyId}", api.GetWorkspaceAccessKeyAction)          // get workspace access key
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanDeleteWorkspaceAccessKey})).Delete("/{keyId}", api.DeleteWorkspaceAccessKeyAction) // delete workspace access key
		})

		if conf.IsSaaS() {
			r.With(middleware.Protect(middleware.Config{Perm: module.CanGetWorkspaceBilling})).Get("/billing", api.GetBillingStatusAction)                               // get billing status
			r.With(middleware.Protect(middleware.Config{Perm: module.CanGetWorkspaceBilling})).Get("/billing/usage", api.GetBillingUsageAction)                          // get billing usage
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanUpdateWorkspaceBilling})).Post("/billing/checkout", api.CreateBillingCheckoutAction) // start Stripe checkout
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanUpdateWorkspaceBilling})).Post("/billing/portal", api.CreateBillingPortalAction)     // open Stripe customer portal
		}

		r.With(middleware.Protect(middleware.Config{Perm: module.CanGetWorkspace})).Get("/stats", api.GetWorkspaceStatsAction) // get workspace stats

		r.Route("/audits", func(r chi.Router) {
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanListWorkspaceAudits})).Get("/", api.ListWorkspaceAuditsAction)      // list audit events
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanGetWorkspaceAudit})).Get("/{auditId}", api.GetWorkspaceAuditAction) // get audit event
		})

		r.Route("/documents", func(r chi.Router) {
			r.With(middleware.Protect(middleware.Config{Perm: module.CanListWorkspaceDocuments})).Get("/", api.ListDocumentsAction)                              // list knowledge documents
			r.With(middleware.Protect(middleware.Config{Perm: module.CanQueryWorkspaceDocuments})).Post("/search", api.SearchDocumentsAction)                    // semantic document search
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanCreateWorkspaceDocument})).Post("/", api.UploadDocumentAction)               // upload knowledge document
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanDeleteWorkspaceDocument})).Delete("/{documentId}", api.DeleteDocumentAction) // delete knowledge document
		})
	})

	r.With(middleware.BasicAuth(
		viper.GetString("app.metrics.username"),
		viper.GetString("app.metrics.secret"),
	)).Get(
		"/api/v1/public/_metrics", promhttp.Handler().ServeHTTP, // Prometheus metrics
	)

	dist, err := fs.Sub(Static, "web/dist")
	if err != nil {
		panic(fmt.Sprintf(
			"Error while accessing dist files: %s",
			err.Error(),
		))
	}

	r.Handle("/assets/*", http.StripPrefix("/", http.FileServer(http.FS(dist)))) // static frontend assets
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {                    // SPA fallback — serve index.html
		indexFile, err := dist.Open("index.html")
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		stat, err := indexFile.Stat()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
	})

	return r
}

// Run starts the HTTP server with graceful shutdown support
func RunServer(handler http.Handler) error {
	err := db.InitDB(ReadWriteDatabase(), ReadOnlyDatabase()...)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	err = github.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize github app: %w", err)
	}

	module.RegisterEventListeners()

	err = module.StartBus()
	if err != nil {
		return fmt.Errorf("failed to start nats bus: %w", err)
	}

	defer func() {
		err := db.CloseDB()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error closing database connection")
		}
	}()

	defer module.StopBus()

	timeout := time.Duration(viper.GetInt("app.timeout")) * time.Second

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("app.port"))),
		Handler:           handler,
		MaxHeaderBytes:    230 * 1024,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadTimeout:       timeout + 5*time.Second,
		WriteTimeout:      timeout + 5*time.Second,
	}

	serr := make(chan error, 1)

	go func() {
		log.Info().
			Int("port", viper.GetInt("app.port")).
			Str("edition", conf.Edition()).
			Msg("Starting HTTP server")

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serr:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Info().
			Str("signal", sig.String()).
			Msg("Received shutdown signal")

		shutdownTimeout := 30 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		log.Info().
			Dur("timeout", shutdownTimeout).
			Msg("Gracefully shutting down server")

		err := srv.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		log.Info().Msg("Server shutdown complete")
	}

	return nil
}
