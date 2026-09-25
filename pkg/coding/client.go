// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package coding

import (
	"context"

	"github.com/clivern/swarm"
	"github.com/spf13/viper"
)

// Client runs coding agents through Swarm.
type Client struct {
	dir    string
	model  string
	image  string
	apiKey string
}

// Request is one coding job against a repository.
type Request struct {
	ID       string
	RepoURL  string
	Prompt   string
	Token    string
	Username string
	SSHKey   string
}

// ChangedFile is one path the agent touched.
type ChangedFile struct {
	Path    string
	Status  string
	Content string
}

// Result is the patch Swarm produced.
type Result struct {
	Patch        string
	Summary      string
	TotalTokens  int
	ChangedFiles []ChangedFile
}

// New returns a coding client loaded from app.coding config.
func New() *Client {
	return &Client{
		dir:    viper.GetString("app.coding.dir"),
		model:  viper.GetString("app.coding.model"),
		image:  viper.GetString("app.coding.docker_image"),
		apiKey: viper.GetString("app.ai.api_key"),
	}
}

// Run clones the repository, runs the agent, and returns the patch.
func (c *Client) Run(ctx context.Context, req Request) (*Result, error) {
	out, err := swarm.Run(ctx, swarm.RunRequest{
		WorkDir:          c.dir,
		ID:               req.ID,
		RepoURL:          req.RepoURL,
		Prompt:           req.Prompt,
		PIModel:          c.model,
		OpenRouterAPIKey: c.apiKey,
		DockerImage:      c.image,
		Cleanup:          true,
		GitCloneAuth: swarm.GitCloneAuth{
			Token:             req.Token,
			Username:          req.Username,
			SSHPrivateKeyPath: req.SSHKey,
		},
	})
	if err != nil {
		return nil, err
	}

	files := make([]ChangedFile, len(out.ChangedFiles))
	for i, file := range out.ChangedFiles {
		files[i] = ChangedFile{
			Path: file.Path,
			// D for deleted, A for added, M for modified
			Status:  file.Status,
			Content: file.Content,
		}
	}

	return &Result{
		Patch:        out.Patch,
		Summary:      out.Summary,
		TotalTokens:  out.TotalTokens,
		ChangedFiles: files,
	}, nil
}
