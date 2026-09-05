// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/github"
	"github.com/clivern/ziee/pkg/github/webhook"

	"github.com/rs/zerolog/log"
)

var ErrInstallationNotFound = errors.New("installation not found")

// InstallationResponse is a GitHub App installation shaped for API responses.
type InstallationResponse struct {
	Id                  db.Id  `json:"id"`
	GitHubId            int64  `json:"githubId"`
	GitHubUserId        string `json:"githubUserId"`
	AccountId           int64  `json:"accountId"`
	AccountLogin        string `json:"accountLogin"`
	AccountType         string `json:"accountType"`
	WorkspaceId         db.Id  `json:"workspaceId"`
	Status              string `json:"status"`
	RepositorySelection string `json:"repositorySelection"`
	HTMLURL             string `json:"htmlUrl"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}

// AttachInstallationRequest is the body for attaching an installation to a workspace.
type AttachInstallationRequest struct {
	WorkspaceId string `json:"workspaceId" validate:"required" label:"Workspace"`
}

// Installation is the module for GitHub App installations.
type Installation struct {
	InstallationRepository db.GitHubInstallationRepository
	RepoRepository         db.WorkspaceGitHubRepoRepository
}

// NewInstallation creates an installation module with the given repositories.
func NewInstallation(installations db.GitHubInstallationRepository, repos db.WorkspaceGitHubRepoRepository) *Installation {
	return &Installation{
		InstallationRepository: installations,
		RepoRepository:         repos,
	}
}

// HandleWebhook persists GitHub App installation and repository access changes.
func (i *Installation) HandleWebhook(event string, body []byte) error {
	switch event {
	case "installation":
		return i.HandleInstallation(body)
	case "installation_repositories":
		return i.HandleInstallationRepositories(body)
	default:
		return nil
	}
}

// HandleInstallation handles a GitHub App installation webhook.
func (i *Installation) HandleInstallation(body []byte) error {
	var payload webhook.InstallationEvent
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return fmt.Errorf("decode installation webhook: %w", err)
	}

	if payload.Action == "deleted" {
		err = i.RepoRepository.DeleteByInstallationId(payload.Installation.ID)
		if err != nil {
			return fmt.Errorf("delete installation repos: %w", err)
		}

		err = i.InstallationRepository.DeleteByGitHubId(payload.Installation.ID)
		if err != nil {
			return fmt.Errorf("delete installation: %w", err)
		}

		log.Info().
			Int64("githubId", payload.Installation.ID).
			Msg("GitHub installation deleted")

		return nil
	}

	meta, err := json.Marshal(payload.Installation)
	if err != nil {
		return fmt.Errorf("encode installation meta: %w", err)
	}
	raw := string(meta)

	installation := &db.GitHubInstallation{
		GitHubId:            payload.Installation.ID,
		GitHubUserId:        strconv.FormatInt(payload.Sender.ID, 10),
		AccountId:           payload.Installation.Account.ID,
		AccountLogin:        payload.Installation.Account.Login,
		AccountType:         payload.Installation.Account.Type,
		RepositorySelection: payload.Installation.RepositorySelection,
		HTMLURL:             payload.Installation.HTMLURL,
		Meta:                &raw,
	}

	err = i.InstallationRepository.Upsert(installation)
	if err != nil {
		return fmt.Errorf("upsert installation: %w", err)
	}

	log.Info().
		Str("id", installation.Id.String()).
		Int64("githubId", installation.GitHubId).
		Str("githubUserId", installation.GitHubUserId).
		Str("status", installation.Status).
		Msg("GitHub installation upserted")

	return nil
}

// HandleInstallationRepositories handles a GitHub App installation_repositories webhook.
func (i *Installation) HandleInstallationRepositories(body []byte) error {
	var payload webhook.InstallationRepositoriesEvent
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return fmt.Errorf("decode installation repositories webhook: %w", err)
	}

	item, err := i.InstallationRepository.GetByGitHubId(payload.Installation.ID)
	if err != nil {
		return fmt.Errorf("get installation: %w", err)
	}
	if item == nil || item.WorkspaceId == "" {
		return nil
	}

	for _, repo := range payload.RepositoriesAdded {
		err = i.StoreRepository(item.WorkspaceId, payload.Installation.ID, repo)
		if err != nil {
			return err
		}
	}

	for _, repo := range payload.RepositoriesRemoved {
		err = i.RepoRepository.DeleteByGitHubId(repo.ID)
		if err != nil {
			return fmt.Errorf("delete installation repo: %w", err)
		}
	}

	log.Info().
		Int64("githubId", payload.Installation.ID).
		Str("workspaceId", item.WorkspaceId.String()).
		Str("action", payload.Action).
		Int("added", len(payload.RepositoriesAdded)).
		Int("removed", len(payload.RepositoriesRemoved)).
		Msg("GitHub installation repositories updated")

	return nil
}

// StoreRepository stores a GitHub repository in the database.
func (i *Installation) StoreRepository(workspaceId db.Id, installationId int64, repo webhook.InstallationRepo) error {
	meta, err := json.Marshal(repo)
	if err != nil {
		return fmt.Errorf("encode repo meta: %w", err)
	}
	raw := string(meta)
	owner, _, _ := strings.Cut(repo.FullName, "/")

	err = i.RepoRepository.Upsert(&db.WorkspaceGitHubRepo{
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

// ListPending lists pending GitHub App installations for a GitHub user to attach to a workspace.
func (i *Installation) ListPending(githubUserId string) ([]*InstallationResponse, error) {
	list, err := i.InstallationRepository.ListPendingByGitHubUserId(githubUserId)
	if err != nil {
		return nil, fmt.Errorf("list pending installations: %w", err)
	}

	installations := make([]*InstallationResponse, 0, len(list))
	for _, item := range list {
		installations = append(installations, &InstallationResponse{
			Id:                  item.Id,
			GitHubId:            item.GitHubId,
			GitHubUserId:        item.GitHubUserId,
			AccountId:           item.AccountId,
			AccountLogin:        item.AccountLogin,
			AccountType:         item.AccountType,
			WorkspaceId:         item.WorkspaceId,
			Status:              item.Status,
			RepositorySelection: item.RepositorySelection,
			HTMLURL:             item.HTMLURL,
			CreatedAt:           item.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:           item.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return installations, nil
}

// Attach links a GitHub App installation to a workspace and stores its repos.
func (i *Installation) Attach(ctx context.Context, id, workspaceId db.Id, githubUserId string) error {
	item, err := i.InstallationRepository.GetById(id)
	if err != nil {
		return fmt.Errorf("get installation: %w", err)
	}
	if item == nil || item.GitHubUserId != githubUserId {
		return ErrInstallationNotFound
	}

	repos, err := github.Get().Repositories(ctx, item.GitHubId)
	if err != nil {
		return fmt.Errorf("list installation repos: %w", err)
	}

	for _, repo := range repos {
		err = i.StoreRepository(workspaceId, item.GitHubId, webhook.InstallationRepo{
			ID:       repo.ID,
			NodeID:   repo.NodeID,
			Name:     repo.Name,
			FullName: repo.FullName,
			Private:  repo.Private,
		})
		if err != nil {
			return err
		}
	}

	err = i.InstallationRepository.Attach(id, workspaceId)
	if err != nil {
		return fmt.Errorf("attach installation: %w", err)
	}

	log.Info().
		Str("id", id.String()).
		Str("workspaceId", workspaceId.String()).
		Int("repos", len(repos)).
		Msg("GitHub installation attached")

	return nil
}
