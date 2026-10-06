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
	DefaultMinPort  = 20000
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
func (r *Runner) Start(ctx context.Context, repositoryId db.Id, remoteId string, ttl time.Duration) (*db.Sandbox, error) {
	repo, err := r.repos.GetById(repositoryId)
	if err != nil {
		return nil, err
	}

	installToken, err := r.github.GetInstallationToken(ctx, repo.InstallationId)
	if err != nil {
		return nil, err
	}

	runId, err := db.NewId()
	if err != nil {
		return nil, err
	}

	token, err := db.NewId()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(r.config.TempDir, runId.String())
	err = Clone(ctx, repo.FullName, installToken.Token, dir)
	if err != nil {
		return nil, err
	}

	port, err := FreePort(r.config.MinPort)
	if err != nil {
		return nil, err
	}

	err = r.Run(ctx, runId.String(), token.String(), dir, port)
	if err != nil {
		return nil, err
	}

	config, err := json.Marshal(RunConfig{
		Image: r.config.DockerImage,
		Model: r.config.Model,
		Dir:   dir,
	})
	if err != nil {
		return nil, err
	}

	item := &db.Sandbox{
		RepositoryId: repositoryId,
		Config:       string(config),
		Port:         port,
		Token:        token.String(),
		RemoteId:     remoteId,
		RunId:        runId,
		ExpiresAt:    time.Now().UTC().Add(ttl),
	}

	err = r.sandboxes.Create(item)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("sandboxId", item.Id.String()).
		Str("repositoryId", repositoryId.String()).
		Str("remoteId", item.RemoteId).
		Str("runId", item.RunId.String()).
		Int("port", port).
		Msg("Sandbox started")

	return item, nil
}

// Stop removes the sandbox container, its clone and its row.
func (r *Runner) Stop(ctx context.Context, id db.Id) error {
	item, err := r.sandboxes.GetById(id)
	if err != nil {
		return err
	}

	Block(item.RunId.String())

	err = exec.CommandContext(ctx, "docker", "rm", "-f", fmt.Sprintf("%s%s", ContainerPrefix, item.RunId.String())).Run()
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(r.config.TempDir, item.RunId.String()))
	if err != nil {
		return err
	}

	err = r.sandboxes.Delete(id)
	if err != nil {
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

	out, err := exec.CommandContext(
		ctx,
		"docker", "run", "-d",
		"--name", fmt.Sprintf("%s%s", ContainerPrefix, runId),
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

	return nil
}

// Clone shallow clones a GitHub repository into dir using an installation token.
func Clone(ctx context.Context, fullName, token, dir string) error {
	_, err := git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:   fmt.Sprintf("https://github.com/%s.git", fullName),
		Depth: 1,
		Auth: &githttp.BasicAuth{
			Username: "x-access-token",
			Password: token,
		},
	})

	return err
}

// FreePort returns the first free local TCP port starting from minPort.
func FreePort(minPort int) (int, error) {
	if minPort <= 0 {
		minPort = DefaultMinPort
	}

	for port := minPort; port <= 65535; port++ {
		l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			continue
		}
		l.Close()

		return port, nil
	}

	return 0, fmt.Errorf("no free port found from %d", minPort)
}
