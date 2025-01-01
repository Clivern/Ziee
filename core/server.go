// Copyright 2026 Clivern. All rights reserved.
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

	"github.com/actx0/ziee/api"
	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/middleware"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/ai"
	"github.com/actx0/ziee/pkg/qdrant"
	"github.com/actx0/ziee/pkg/storage"
	"github.com/actx0/ziee/service/agent"
	"github.com/actx0/ziee/service/knowledge"
	"github.com/actx0/ziee/task"

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
		r.Get("/api/v1/public/_health", api.HealthAction)                         // liveness probe
		r.Get("/api/v1/public/_ready", api.ReadyAction)                           // readiness probe (db, deps)
		r.Post("/api/v1/public/action/setup", api.SetupAction)                    // initial app setup
		r.Get("/api/v1/public/action/setup/status", api.SetupStatusAction)        // setup completion status
		r.Post("/api/v1/public/action/login", api.LoginAction)                    // user login
		r.Post("/api/v1/public/action/logout", api.LogoutAction)                  // user logout
		r.Post("/api/v1/public/action/forgot-password", api.ForgotPasswordAction) // request password reset email
		r.Post("/api/v1/public/action/reset-password", api.ResetPasswordAction)   // reset password with token
		r.Post("/api/v1/public/action/register", api.RegisterAction)              // user registration
		r.Post("/api/v1/public/action/stripe/webhook", api.StripeWebhookAction)   // Stripe billing webhook
	})
	r.Get("/api/v1/me", api.GetMeAction) // current authenticated user
	r.Group(func(r chi.Router) {         // user profile and workspace invites
		r.Use(middleware.Protect(middleware.Config{Roles: []string{db.UserRoleAdmin, db.UserRoleRegular}}))
		r.Get("/api/v1/action/profile", api.GetProfileAction)                                        // get user profile
		r.Put("/api/v1/action/profile", api.UpdateProfileAction)                                     // update user profile
		r.Get("/api/v1/action/invites", api.ListUserInvitesAction)                                   // list pending workspace invites
		r.Get("/api/v1/action/invite-by-token/{token}", api.GetAuthenticatedUserInviteByTokenAction) // get invite details by token
		r.Post("/api/v1/action/accept-invite/{token}", api.AcceptUserInviteByTokenAction)            // accept workspace invite
		r.Post("/api/v1/action/reject-invite/{token}", api.RejectUserInviteByTokenAction)            // reject workspace invite
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
		r.Use(middleware.TrackWorkspaceAPICall())

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

		r.With(middleware.Protect(middleware.Config{Perm: module.CanGetWorkspaceBilling})).Get("/billing", api.GetBillingStatusAction)                               // get billing status
		r.With(middleware.Protect(middleware.Config{Perm: module.CanGetWorkspaceBilling})).Get("/billing/usage", api.GetBillingUsageAction)                          // get billing usage
		r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanUpdateWorkspaceBilling})).Post("/billing/checkout", api.CreateBillingCheckoutAction) // start Stripe checkout
		r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanUpdateWorkspaceBilling})).Post("/billing/portal", api.CreateBillingPortalAction)     // open Stripe customer portal

		r.With(middleware.Protect(middleware.Config{Perm: module.CanGetWorkspace})).Get("/stats", api.GetWorkspaceStatsAction) // get workspace stats

		r.Route("/audits", func(r chi.Router) {
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanListWorkspaceAudits})).Get("/", api.ListWorkspaceAuditsAction)      // list audit events
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanGetWorkspaceAudit})).Get("/{auditId}", api.GetWorkspaceAuditAction) // get audit event
		})

		r.Route("/prompts", func(r chi.Router) {
			r.With(middleware.Protect(middleware.Config{Perm: module.CanListPrompts})).Get("/", api.ListPromptsAction)                            // list prompts
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanCreatePrompt})).Post("/", api.CreatePromptAction)             // create prompt
			r.With(middleware.Protect(middleware.Config{Perm: module.CanGetPrompt})).Get("/{promptId}", api.GetPromptAction)                      // get prompt
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanDeletePrompt})).Delete("/{promptId}", api.DeletePromptAction) // delete prompt

			r.Route("/{promptId}/versions", func(r chi.Router) {
				r.With(middleware.Protect(middleware.Config{Perm: module.CanListPrompts})).Get("/", api.ListPromptVersionsAction)                                   // list prompt versions
				r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanCreatePrompt})).Post("/", api.CreatePromptVersionAction)                    // create prompt version
				r.With(middleware.Protect(middleware.Config{Perm: module.CanGetPrompt})).Get("/{promptVersionId}", api.GetPromptVersionAction)                      // get prompt version
				r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanUpdatePrompt})).Put("/{promptVersionId}", api.UpdatePromptVersionAction)    // update prompt version
				r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanDeletePrompt})).Delete("/{promptVersionId}", api.DeletePromptVersionAction) // delete prompt version
			})
		})

		r.With(middleware.Protect(middleware.Config{Perm: module.CanGetPrompt})).Get("/promptsByName/{promptName}", api.GetPromptByNameAction) // get prompt by name

		r.Route("/documents", func(r chi.Router) {
			r.With(middleware.Protect(middleware.Config{Perm: module.CanListWorkspaceDocuments})).Get("/", api.ListDocumentsAction)                              // list knowledge documents
			r.With(middleware.Protect(middleware.Config{Perm: module.CanQueryWorkspaceDocuments})).Post("/search", api.SearchDocumentsAction)                    // semantic document search
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanCreateWorkspaceDocument})).Post("/", api.UploadDocumentAction)               // upload knowledge document
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanDeleteWorkspaceDocument})).Delete("/{documentId}", api.DeleteDocumentAction) // delete knowledge document
		})

		r.Route("/agents", func(r chi.Router) {
			r.With(middleware.Protect(middleware.Config{Perm: module.CanListAgents})).Get("/", api.ListAgentsAction)                        // list agents
			r.With(middleware.Protect(middleware.Config{Perm: module.CanCreateAgent})).Post("/", api.CreateAgentAction)                     // create agent
			r.With(middleware.Protect(middleware.Config{Perm: module.CanGetAgent})).Get("/{agentId}", api.GetAgentAction)                   // get agent
			r.With(middleware.Protect(middleware.Config{User: true, Perm: module.CanUpdateAgent})).Put("/{agentId}", api.UpdateAgentAction) // update agent
			r.With(middleware.Protect(middleware.Config{Perm: module.CanDeleteAgent})).Delete("/{agentId}", api.DeleteAgentAction)          // delete agent

			r.Route("/{agentId}/sessions", func(r chi.Router) {
				r.With(middleware.Protect(middleware.Config{Perm: module.CanListAgentSessions})).Get("/", api.ListAgentSessionsAction)                       // list agent sessions
				r.With(middleware.Protect(middleware.Config{Perm: module.CanGetAgentSession})).Get("/by-labels", api.GetAgentSessionByLabelsAction)          // get session by labels
				r.With(middleware.Protect(middleware.Config{Perm: module.CanGetAgentSession})).Get("/{sessionId}", api.GetAgentSessionAction)                // get agent session
				r.With(middleware.Protect(middleware.Config{Perm: module.CanCreateAgentSession})).Post("/", api.CreateAgentSessionAction)                    // create agent session
				r.With(middleware.Protect(middleware.Config{Perm: module.CanUpdateAgentSession})).Put("/by-labels", api.UpdateAgentSessionByLabelsAction)    // update session by labels
				r.With(middleware.Protect(middleware.Config{Perm: module.CanDeleteAgentSession})).Delete("/by-labels", api.DeleteAgentSessionByLabelsAction) // delete session by labels

				r.Route("/{sessionId}/messages", func(r chi.Router) {
					r.With(middleware.Protect(middleware.Config{Perm: module.CanListSessionMessages})).Get("/", api.ListSessionMessagesAction)                 // list session messages
					r.With(middleware.Protect(middleware.Config{Perm: module.CanCreateSessionMessage})).Post("/", api.CreateSessionMessageAction)              // create session message
					r.With(middleware.Protect(middleware.Config{Perm: module.CanCreateSessionMessage})).Post("/batch", api.BatchCreateSessionMessagesAction)   // batch create session messages
					r.With(middleware.Protect(middleware.Config{Perm: module.CanDeleteSessionMessage})).Delete("/batch", api.BatchDeleteSessionMessagesAction) // batch delete session messages
					r.With(middleware.Protect(middleware.Config{Perm: module.CanGetSessionMessage})).Get("/{messageId}", api.GetSessionMessageAction)          // get session message
					r.With(middleware.Protect(middleware.Config{Perm: module.CanUpdateSessionMessage})).Put("/{messageId}", api.UpdateSessionMessageAction)    // update session message
					r.With(middleware.Protect(middleware.Config{Perm: module.CanDeleteSessionMessage})).Delete("/{messageId}", api.DeleteSessionMessageAction) // delete session message
				})

				r.Route("/{sessionId}/memories", func(r chi.Router) {
					r.With(middleware.Protect(middleware.Config{Perm: module.CanListSessionMemories})).Get("/", api.ListSessionMemoriesAction)                // list session memories
					r.With(middleware.Protect(middleware.Config{Perm: module.CanCreateSessionMemory})).Post("/", api.CreateSessionMemoryAction)               // create session memory
					r.With(middleware.Protect(middleware.Config{Perm: module.CanCreateSessionMemory})).Post("/batch", api.BatchCreateSessionMemoriesAction)   // batch create session memories
					r.With(middleware.Protect(middleware.Config{Perm: module.CanDeleteSessionMemory})).Delete("/batch", api.BatchDeleteSessionMemoriesAction) // batch delete session memories
					r.With(middleware.Protect(middleware.Config{Perm: module.CanGetSessionMemory})).Get("/{memoryId}", api.GetSessionMemoryAction)            // get session memory
					r.With(middleware.Protect(middleware.Config{Perm: module.CanUpdateSessionMemory})).Put("/{memoryId}", api.UpdateSessionMemoryAction)      // update session memory
					r.With(middleware.Protect(middleware.Config{Perm: module.CanDeleteSessionMemory})).Delete("/{memoryId}", api.DeleteSessionMemoryAction)   // delete session memory
				})
			})
		})
	})

	r.Group(func(r chi.Router) { // demo weather tools
		r.Use(middleware.Protect(middleware.Config{Roles: []string{db.UserRoleAdmin, db.UserRoleRegular}}))
		r.Post("/api/v1/tool/getCityInformation", api.GetCityInformationAction)           // get city info by name
		r.Post("/api/v1/tool/getWeatherByCity", api.GetWeatherByCityAction)               // get weather by city name
		r.Post("/api/v1/tool/getWeatherByCoordinates", api.GetWeatherByCoordinatesAction) // get weather by coordinates
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

	module.RegisterEventListeners()

	asyr, err := module.Start(db.NewAsyncTaskRepository(db.GetDB(true)))
	if err != nil {
		return fmt.Errorf("failed to start async worker pool: %w", err)
	}

	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to initialize document storage: %w", err)
	}

	vdb, err := qdrant.New()
	if err != nil {
		return fmt.Errorf("failed to initialize qdrant: %w", err)
	}

	defer func() {
		err := vdb.Close()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error closing qdrant client")
		}
	}()

	ksvc := knowledge.New(knowledge.Dependencies{
		Documents: db.NewWorkspaceDocumentRepository(
			db.GetDB(true),
		),
		Embed:         ai.NewEmbedClient(),
		Vectors:       vdb,
		Store:         store,
		Usage:         db.NewUsageRepository(db.GetDB(false)),
		Subscriptions: db.NewSubscriptionRepository(db.GetDB(false)),
	})

	asvc := agent.New(agent.Dependencies{
		Embed:    ai.NewEmbedClient(),
		Vectors:  vdb,
		LiteLLM:  ai.NewLiteClient(),
		LargeLLM: ai.NewLargeClient(),
		Agents:   db.NewAgentRepository(db.GetDB(true)),
		Sessions: db.NewAgentSessionRepository(db.GetDB(true)),
		Messages: db.NewSessionMessageRepository(db.GetDB(true)),
		Memories: db.NewSessionMemoryRepository(db.GetDB(true)),
	})

	task.Register(asyr, task.Dependencies{
		Knowledge: ksvc,
		Agent:     asvc,
	})

	defer func() {
		err := db.CloseDB()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error closing database connection")
		}
	}()

	defer func() {
		err := asyr.Stop(30 * time.Second)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error stopping async worker pool")
		}
	}()

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
