// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"errors"
	"fmt"
	"time"

	"github.com/clivern/ziee/db"

	"github.com/rs/zerolog/log"
)

var (
	ErrSandboxNotFound     = errors.New("sandbox not found")
	ErrFailedCreateSandbox = errors.New("failed create sandbox")
	ErrFailedListSandboxes = errors.New("failed list sandboxes")
	ErrFailedGetSandbox    = errors.New("failed get sandbox")
	ErrFailedDeleteSandbox = errors.New("failed delete sandbox")
	ErrFailedTouchSandbox  = errors.New("failed touch sandbox")
)

// Sandbox is the module for repository sandbox CRUD.
type Sandbox struct {
	SandboxRepository db.SandboxRepository
	RepoRepository    db.RepositoriesRepository
}

// NewSandbox creates a sandbox module with the given repositories.
func NewSandbox(sandboxes db.SandboxRepository, repos db.RepositoriesRepository) *Sandbox {
	return &Sandbox{
		SandboxRepository: sandboxes,
		RepoRepository:    repos,
	}
}

// SandboxResponse is a sandbox shaped for API responses.
type SandboxResponse struct {
	Id             db.Id  `json:"id"`
	RepositoryId   db.Id  `json:"repositoryId"`
	Config         string `json:"config"`
	Usage          string `json:"usage"`
	Port           int    `json:"port"`
	Status         string `json:"status"`
	Token          string `json:"token,omitempty"`
	RemoteId       string `json:"remoteId,omitempty"`
	ExpiresAt      string `json:"expiresAt"`
	LastActivityAt string `json:"lastActivityAt"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// ListSandboxesResult is what you get back when listing sandboxes.
type ListSandboxesResult struct {
	Sandboxes []*SandboxResponse
	Total     int
}

// CreateSandbox creates a sandbox for a repository.
func (s *Sandbox) CreateSandbox(repositoryId db.Id, config, usage, token, remoteId string, port int, expiresAt time.Time) (*SandboxResponse, error) {
	repo, err := s.RepoRepository.GetById(repositoryId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateSandbox, err)
	}
	if repo == nil {
		return nil, fmt.Errorf("%w: repository not found", ErrFailedCreateSandbox)
	}

	item := &db.Sandbox{
		RepositoryId: repositoryId,
		Config:       config,
		Usage:        usage,
		Port:         port,
		Token:        token,
		RemoteId:     remoteId,
		ExpiresAt:    expiresAt,
	}

	err = s.SandboxRepository.Create(item)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateSandbox, err)
	}

	log.Info().
		Str("sandboxId", item.Id.String()).
		Str("repositoryId", repositoryId.String()).
		Str("remoteId", remoteId).
		Msg("Sandbox created")

	return &SandboxResponse{
		Id:             item.Id,
		RepositoryId:   item.RepositoryId,
		Config:         item.Config,
		Usage:          item.Usage,
		Port:           item.Port,
		Status:         item.Status,
		Token:          item.Token,
		RemoteId:       item.RemoteId,
		ExpiresAt:      item.ExpiresAt.UTC().Format(time.RFC3339),
		LastActivityAt: item.LastActivityAt.UTC().Format(time.RFC3339),
		CreatedAt:      item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      item.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetSandbox returns a sandbox by id.
func (s *Sandbox) GetSandbox(id db.Id) (*SandboxResponse, error) {
	item, err := s.SandboxRepository.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetSandbox, err)
	}
	if item == nil {
		return nil, ErrSandboxNotFound
	}

	return &SandboxResponse{
		Id:             item.Id,
		RepositoryId:   item.RepositoryId,
		Config:         item.Config,
		Usage:          item.Usage,
		Port:           item.Port,
		Status:         item.Status,
		Token:          item.Token,
		RemoteId:       item.RemoteId,
		ExpiresAt:      item.ExpiresAt.UTC().Format(time.RFC3339),
		LastActivityAt: item.LastActivityAt.UTC().Format(time.RFC3339),
		CreatedAt:      item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      item.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetSandboxByToken returns a sandbox by token.
func (s *Sandbox) GetSandboxByToken(token string) (*SandboxResponse, error) {
	item, err := s.SandboxRepository.GetByToken(token)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetSandbox, err)
	}
	if item == nil {
		return nil, ErrSandboxNotFound
	}

	return &SandboxResponse{
		Id:             item.Id,
		RepositoryId:   item.RepositoryId,
		Config:         item.Config,
		Usage:          item.Usage,
		Port:           item.Port,
		Status:         item.Status,
		Token:          item.Token,
		RemoteId:       item.RemoteId,
		ExpiresAt:      item.ExpiresAt.UTC().Format(time.RFC3339),
		LastActivityAt: item.LastActivityAt.UTC().Format(time.RFC3339),
		CreatedAt:      item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      item.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// ListSandboxes returns sandboxes for a repository.
func (s *Sandbox) ListSandboxes(repositoryId db.Id) (*ListSandboxesResult, error) {
	items, err := s.SandboxRepository.ListByRepositoryId(repositoryId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListSandboxes, err)
	}

	list := make([]*SandboxResponse, 0, len(items))
	for _, item := range items {
		list = append(list, &SandboxResponse{
			Id:             item.Id,
			RepositoryId:   item.RepositoryId,
			Config:         item.Config,
			Usage:          item.Usage,
			Port:           item.Port,
			Status:         item.Status,
			Token:          item.Token,
			RemoteId:       item.RemoteId,
			ExpiresAt:      item.ExpiresAt.UTC().Format(time.RFC3339),
			LastActivityAt: item.LastActivityAt.UTC().Format(time.RFC3339),
			CreatedAt:      item.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:      item.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	log.Info().
		Str("repositoryId", repositoryId.String()).
		Int("count", len(list)).
		Msg("Sandboxes listed")

	return &ListSandboxesResult{Sandboxes: list, Total: len(list)}, nil
}

// TouchSandbox updates last activity for a sandbox.
func (s *Sandbox) TouchSandbox(id db.Id) error {
	item, err := s.SandboxRepository.GetById(id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedTouchSandbox, err)
	}
	if item == nil {
		return ErrSandboxNotFound
	}

	err = s.SandboxRepository.Touch(id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedTouchSandbox, err)
	}

	log.Info().
		Str("sandboxId", id.String()).
		Msg("Sandbox touched")

	return nil
}

// DeleteSandbox removes a sandbox.
func (s *Sandbox) DeleteSandbox(id db.Id) error {
	item, err := s.SandboxRepository.GetById(id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteSandbox, err)
	}
	if item == nil {
		return ErrSandboxNotFound
	}

	err = s.SandboxRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteSandbox, err)
	}

	log.Info().
		Str("sandboxId", id.String()).
		Str("repositoryId", item.RepositoryId.String()).
		Msg("Sandbox deleted")

	return nil
}

// DeleteExpired removes expired sandboxes.
func (s *Sandbox) DeleteExpired() (int64, error) {
	count, err := s.SandboxRepository.DeleteExpired()
	if err != nil {
		return 0, err
	}

	if count > 0 {
		log.Info().
			Int64("count", count).
			Msg("Expired sandboxes deleted")
	}

	return count, nil
}
