// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/github/app"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// Repository is the module for GitHub repositories.
type Repository struct {
	RepoRepository     db.RepositoriesRepository
	MetaRepository     db.RepositoryMetaRepository
	SpamUserRepository db.RepositorySpamUserRepository
}

// NewRepository creates a repository module with the given stores.
func NewRepository(
	repos db.RepositoriesRepository,
	meta db.RepositoryMetaRepository,
	spam db.RepositorySpamUserRepository,
) *Repository {
	return &Repository{
		RepoRepository:     repos,
		MetaRepository:     meta,
		SpamUserRepository: spam,
	}
}

// Import stores repositories and enqueues a bootstrap task for each.
func (r *Repository) Import(workspaceId db.Id, installationId int64, repositories []app.Repository) error {
	for _, repository := range repositories {
		err := r.Store(workspaceId, installationId, repository)
		if err != nil {
			return err
		}

		err = EnqueueTask(db.AsyncTaskTypeRepoBootstrap, map[string]string{
			"workspaceId":    workspaceId.String(),
			"installationId": strconv.FormatInt(installationId, 10),
			"githubId":       strconv.FormatInt(repository.ID, 10),
			"fullName":       repository.FullName,
			"name":           repository.Name,
		}, workspaceId)
		if err != nil {
			log.Error().
				Err(err).
				Int64("installationId", installationId).
				Str("repository", repository.FullName).
				Msg("Failed to enqueue repository bootstrap task")
		}
	}

	return nil
}

// Update adds and removes repos for an attached installation.
func (r *Repository) Update(workspaceId db.Id, installationId int64, added []app.Repository, removed []int64) error {
	err := r.Import(workspaceId, installationId, added)
	if err != nil {
		return err
	}

	for _, repositoryId := range removed {
		err = r.RepoRepository.DeleteByGitHubId(repositoryId)
		if err != nil {
			return fmt.Errorf("delete installation repo: %w", err)
		}
	}

	log.Info().
		Int64("installationId", installationId).
		Str("workspaceId", workspaceId.String()).
		Int("added", len(added)).
		Int("removed", len(removed)).
		Msg("GitHub installation repositories updated")

	return nil
}

// Store stores a GitHub repository in the database.
func (r *Repository) Store(workspaceId db.Id, installationId int64, repo app.Repository) error {
	meta, err := json.Marshal(repo)
	if err != nil {
		return fmt.Errorf("encode repo meta: %w", err)
	}

	raw := string(meta)
	owner, _, _ := strings.Cut(repo.FullName, "/")

	err = r.RepoRepository.Upsert(&db.Repository{
		WorkspaceId:    workspaceId,
		InstallationId: installationId,
		GitHubId:       repo.ID,
		NodeId:         repo.NodeID,
		Owner:          owner,
		Name:           repo.Name,
		FullName:       repo.FullName,
		Private:        repo.Private,
		Meta:           &raw,
	})
	if err != nil {
		return fmt.Errorf("store installation repo: %w", err)
	}

	return nil
}

// UpsertMeta stores a repositories_meta value by GitHub repository id.
func (r *Repository) UpsertMeta(githubRepoId int64, key, value string) error {
	repo, err := r.RepoRepository.GetByGitHubId(githubRepoId)
	if err != nil {
		return fmt.Errorf("get repo: %w", err)
	}

	err = r.MetaRepository.Upsert(repo.Id, key, value)
	if err != nil {
		return fmt.Errorf("upsert repo meta: %w", err)
	}

	return nil
}

// GetMeta returns a repositories_meta value by GitHub repository id.
func (r *Repository) GetMeta(githubRepoId int64, key string) (string, error) {
	repo, err := r.RepoRepository.GetByGitHubId(githubRepoId)
	if err != nil {
		return "", fmt.Errorf("get repo: %w", err)
	}

	meta, err := r.MetaRepository.Get(repo.Id, key)
	if err != nil {
		return "", fmt.Errorf("get repo meta: %w", err)
	}

	return lo.FromPtr(meta).Value, nil
}

// SetConfigPath stores the `.ziee.yml` path by GitHub repository id.
func (r *Repository) SetConfigPath(githubRepoId int64, path string) error {
	repo, err := r.RepoRepository.GetByGitHubId(githubRepoId)
	if err != nil {
		return fmt.Errorf("get repo: %w", err)
	}

	repo.ConfigPath = &path
	err = r.RepoRepository.Update(repo)
	if err != nil {
		return fmt.Errorf("set config path: %w", err)
	}

	return nil
}

// WorkspaceId returns the Ziee workspace for a GitHub repository id.
func (r *Repository) WorkspaceId(githubRepoId int64) db.Id {
	repo, err := r.RepoRepository.GetByGitHubId(githubRepoId)
	if err != nil {
		log.Error().
			Err(err).
			Int64("githubRepoId", githubRepoId).
			Msg("Failed to load repository workspace")
	}

	return repo.WorkspaceId
}

// GetConfigPath returns the `.ziee.yml` path by GitHub repository id.
func (r *Repository) GetConfigPath(githubRepoId int64) (string, error) {
	repo, err := r.RepoRepository.GetByGitHubId(githubRepoId)
	if err != nil {
		return "", fmt.Errorf("get repo: %w", err)
	}

	return lo.FromPtr(repo.ConfigPath), nil
}

// IsAuthorBlocked reports whether the GitHub user is on the repository spam blocklist.
func (r *Repository) IsAuthorBlocked(githubRepoId, githubUserId int64) bool {
	repo, err := r.RepoRepository.GetByGitHubId(githubRepoId)
	if err != nil {
		log.Error().
			Err(err).
			Int64("githubRepoId", githubRepoId).
			Int64("githubUserId", githubUserId).
			Msg("Failed to load repository for blocklist check")

		return false
	}

	item, err := r.SpamUserRepository.GetByGitHubId(repo.Id, githubUserId)
	if err != nil {
		log.Error().
			Err(err).
			Int64("githubRepoId", githubRepoId).
			Int64("githubUserId", githubUserId).
			Msg("Failed to load author blocklist")

		return false
	}

	return item != nil
}

// BlockAuthor adds the GitHub user to the repository spam blocklist.
func (r *Repository) BlockAuthor(githubRepoId int64, username string, githubUserId int64) {
	repo, err := r.RepoRepository.GetByGitHubId(githubRepoId)
	if err != nil {
		log.Error().
			Err(err).
			Int64("githubRepoId", githubRepoId).
			Str("author", username).
			Msg("Failed to load repository for blocklist")
	}

	err = r.SpamUserRepository.Upsert(&db.RepositorySpamUser{
		RepositoryId: repo.Id,
		GitHubId:     &githubUserId,
		Username:     &username,
	})
	if err != nil {
		log.Error().
			Err(err).
			Int64("githubRepoId", githubRepoId).
			Str("author", username).
			Msg("Failed to block author")
	}
}
