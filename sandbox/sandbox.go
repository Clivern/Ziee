// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/github/app"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	RPCPort         = 8765
	ContainerPrefix = "ziee-"
)

// Runner clones a repository and runs the sandbox container on it.
type Runner struct {
	config    Config
	github    *app.App
	repos     db.RepositoriesRepository
	sandboxes db.SandboxRepository
}

// RunConfig is the sandbox config stored with the sandbox row.
type RunConfig struct {
	Image string `json:"image"`
	Model string `json:"model"`
	Dir   string `json:"dir"`
}

// StartRequest is what you pass when starting a sandbox.
type StartRequest struct {
	RepositoryId db.Id
	Owner        string
	RemoteId     string
	TTL          time.Duration
}

// NewRunner returns a sandbox runner.
func NewRunner(github *app.App, repos db.RepositoriesRepository, sandboxes db.SandboxRepository) *Runner {
	return &Runner{
		config:    GetConfig(),
		github:    github,
		repos:     repos,
		sandboxes: sandboxes,
	}
}

// Start clones the repository and runs a sandbox container on it.
func (r *Runner) Start(ctx context.Context, req *StartRequest) (*db.Sandbox, error) {
	log.Info().
		Str("repositoryId", req.RepositoryId.String()).
		Str("owner", req.Owner).
		Str("remoteId", req.RemoteId).
		Dur("ttl", req.TTL).
		Msg("Starting sandbox")

	repo, err := r.repos.GetById(req.RepositoryId)
	if err != nil {
		log.Error().
			Err(err).
			Str("repositoryId", req.RepositoryId.String()).
			Msg("Failed to get sandbox repository")
		return nil, err
	}

	installToken, err := r.github.GetInstallationToken(ctx, repo.InstallationId)
	if err != nil {
		log.Error().
			Err(err).
			Str("repositoryId", req.RepositoryId.String()).
			Int64("installationId", repo.InstallationId).
			Msg("Failed to get installation token for sandbox")
		return nil, err
	}

	runId, err := db.NewId()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate sandbox run id")
		return nil, err
	}

	token, err := db.NewId()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate sandbox token")
		return nil, err
	}

	dir := filepath.Join(r.config.TempDir, runId.String())

	// Remove the clone and container if Start fails after this point.
	started := false
	defer func() {
		if !started {
			r.Cleanup(runId.String(), dir)
		}
	}()

	err = Clone(ctx, repo.FullName, installToken.Token, dir)
	if err != nil {
		log.Error().
			Err(err).
			Str("repository", repo.FullName).
			Str("runId", runId.String()).
			Str("dir", dir).
			Msg("Failed to clone sandbox repository")
		return nil, err
	}

	port, err := r.FreePort(r.config.MinPort)
	if err != nil {
		log.Error().
			Err(err).
			Str("runId", runId.String()).
			Int("minPort", r.config.MinPort).
			Msg("Failed to find a free port for sandbox")
		return nil, err
	}

	err = r.Run(ctx, runId.String(), token.String(), dir, port)
	if err != nil {
		log.Error().
			Err(err).
			Str("runId", runId.String()).
			Int("port", port).
			Msg("Failed to run sandbox container")
		return nil, err
	}

	config, err := json.Marshal(RunConfig{
		Image: r.config.DockerImage,
		Model: r.config.Model,
		Dir:   dir,
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("runId", runId.String()).
			Msg("Failed to encode sandbox config")
		return nil, err
	}

	item := &db.Sandbox{
		RepositoryId: req.RepositoryId,
		Owner:        req.Owner,
		Config:       string(config),
		Port:         port,
		Token:        token.String(),
		RemoteId:     req.RemoteId,
		RunId:        runId,
		ExpiresAt:    time.Now().UTC().Add(req.TTL),
	}

	err = r.sandboxes.Create(item)
	if err != nil {
		log.Error().
			Err(err).
			Str("repositoryId", req.RepositoryId.String()).
			Str("runId", runId.String()).
			Msg("Failed to store sandbox")
		return nil, err
	}

	log.Info().
		Str("sandboxId", item.Id.String()).
		Str("repositoryId", req.RepositoryId.String()).
		Str("owner", item.Owner).
		Str("remoteId", item.RemoteId).
		Str("runId", item.RunId.String()).
		Int("port", port).
		Msg("Sandbox started")

	started = true

	return item, nil
}

// Cleanup removes a sandbox container and clone left by a failed Start.
func (r *Runner) Cleanup(runId, dir string) {
	container := fmt.Sprintf("%s%s", ContainerPrefix, runId)

	out, err := exec.Command("docker", "rm", "-f", container).CombinedOutput()
	if err != nil && !strings.Contains(string(out), "No such container") {
		log.Error().
			Err(err).
			Str("runId", runId).
			Str("container", container).
			Str("output", strings.TrimSpace(string(out))).
			Msg("Failed to clean up sandbox container")
	}

	err = os.RemoveAll(dir)
	if err != nil {
		log.Error().
			Err(err).
			Str("runId", runId).
			Str("dir", dir).
			Msg("Failed to clean up sandbox directory")
	}
}

// Stop removes the sandbox container, its clone and its row.
func (r *Runner) Stop(ctx context.Context, id db.Id) error {
	log.Info().
		Str("sandboxId", id.String()).
		Msg("Stopping sandbox")

	item, err := r.sandboxes.GetById(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("sandboxId", id.String()).
			Msg("Failed to get sandbox")
		return err
	}

	Block(item.RunId.String())

	container := fmt.Sprintf("%s%s", ContainerPrefix, item.RunId.String())
	out, err := exec.CommandContext(ctx, "docker", "rm", "-f", container).CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("sandboxId", id.String()).
			Str("container", container).
			Str("output", strings.TrimSpace(string(out))).
			Msg("Failed to remove sandbox container")
		return err
	}

	log.Debug().
		Str("sandboxId", id.String()).
		Str("container", container).
		Msg("Sandbox container removed")

	dir := filepath.Join(r.config.TempDir, item.RunId.String())
	err = os.RemoveAll(dir)
	if err != nil {
		log.Error().
			Err(err).
			Str("sandboxId", id.String()).
			Str("dir", dir).
			Msg("Failed to remove sandbox directory")
		return err
	}

	log.Debug().
		Str("sandboxId", id.String()).
		Str("dir", dir).
		Msg("Sandbox directory removed")

	err = r.sandboxes.Delete(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("sandboxId", id.String()).
			Msg("Failed to delete sandbox")
		return err
	}

	log.Info().
		Str("sandboxId", id.String()).
		Str("remoteId", item.RemoteId).
		Str("runId", item.RunId.String()).
		Msg("Sandbox stopped")

	return nil
}

// Run starts the sandbox container with the repository mounted at /repo.
func (r *Runner) Run(ctx context.Context, runId, token, dir string, port int) error {
	proxyURL := fmt.Sprintf(
		"%s/api/v1/sandbox/%s/api",
		strings.TrimRight(viper.GetString("app.url"), "/"),
		runId,
	)

	container := fmt.Sprintf("%s%s", ContainerPrefix, runId)

	log.Debug().
		Str("runId", runId).
		Str("container", container).
		Str("image", r.config.DockerImage).
		Str("model", r.config.Model).
		Str("dir", dir).
		Int("port", port).
		Str("proxyURL", proxyURL).
		Msg("Running sandbox container")

	out, err := exec.CommandContext(
		ctx,
		"docker", "run", "-d",
		"--name", container,
		"-p", fmt.Sprintf("127.0.0.1:%d:%d", port, RPCPort),
		"-v", fmt.Sprintf("%s:/repo", dir),
		"-e", fmt.Sprintf("RUN_ID=%s", runId),
		"-e", fmt.Sprintf("RPC_API_KEY=%s", token),
		"-e", fmt.Sprintf("PROXY_URL=%s", proxyURL),
		"-e", fmt.Sprintf("PI_MODEL=%s", r.config.Model),
		r.config.DockerImage,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker run: %w: %s", err, out)
	}

	log.Debug().
		Str("runId", runId).
		Str("container", container).
		Str("containerId", strings.TrimSpace(string(out))).
		Msg("Sandbox container running")

	return nil
}

// Clone shallow clones a GitHub repository into dir using an installation token.
func Clone(ctx context.Context, fullName, token, dir string) error {
	start := time.Now()

	log.Debug().
		Str("repository", fullName).
		Str("dir", dir).
		Msg("Cloning sandbox repository")

	_, err := git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:   fmt.Sprintf("https://github.com/%s.git", fullName),
		Depth: 1,
		Auth: &githttp.BasicAuth{
			Username: "x-access-token",
			Password: token,
		},
	})
	if err != nil {
		return err
	}

	log.Debug().
		Str("repository", fullName).
		Str("dir", dir).
		Dur("duration", time.Since(start)).
		Msg("Sandbox repository cloned")

	return nil
}

// FreePort returns the first free local TCP port starting from minPort,
// skipping ports held by active or stopped sandboxes.
func (r *Runner) FreePort(minPort int) (int, error) {
	if minPort <= 0 || minPort > 65535 {
		return 0, fmt.Errorf("invalid min port %d", minPort)
	}

	for port := minPort; port <= 65535; port++ {
		l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			continue
		}

		l.Close()

		reserved, err := r.sandboxes.IsPortReserved(port)
		if err != nil {
			return 0, err
		}
		if reserved {
			continue
		}

		log.Debug().
			Int("minPort", minPort).
			Int("port", port).
			Msg("Found free sandbox port")

		return port, nil
	}

	return 0, fmt.Errorf("no free port found from %d", minPort)
}
