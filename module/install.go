// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/clivern/ziee/db"
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

// HandleWebhook persists or deletes a GitHub App installation from a webhook delivery.
func (i *Installation) HandleWebhook(event string, body []byte) error {
	if event != "installation" {
		return nil
	}

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

// AttachInstallationRequest is the body for attaching an installation to a workspace.
type AttachInstallationRequest struct {
	WorkspaceId string `json:"workspaceId" validate:"required" label:"Workspace"`
}

// Attach links a GitHub App installation to a workspace.
func (i *Installation) Attach(id, workspaceId db.Id, githubUserId string) error {
	item, err := i.InstallationRepository.GetById(id)
	if err != nil {
		return fmt.Errorf("get installation: %w", err)
	}
	if item == nil || item.GitHubUserId != githubUserId {
		return ErrInstallationNotFound
	}

	err = i.InstallationRepository.Attach(id, workspaceId)
	if err != nil {
		return fmt.Errorf("attach installation: %w", err)
	}

	log.Info().
		Str("id", id.String()).
		Str("workspaceId", workspaceId.String()).
		Msg("GitHub installation attached")

	return nil
}
