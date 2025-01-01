// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/util"

	"github.com/samber/lo"
)

var (
	ErrPromptNotFound                   = errors.New("prompt not found")
	ErrPromptVersionNotFound            = errors.New("prompt version not found")
	ErrFailedCreatePromptVersion        = errors.New("failed create prompt version")
	ErrFailedListPrompts                = errors.New("failed list prompts")
	ErrFailedListPromptVersions         = errors.New("failed list prompt versions")
	ErrFailedGetPrompt                  = errors.New("failed get prompt")
	ErrFailedGetPromptVersion           = errors.New("failed get prompt version")
	ErrFailedUpdatePromptVersion        = errors.New("failed update prompt version")
	ErrFailedSetPromptVersionProduction = errors.New("failed set prompt version production")
	ErrFailedDeletePromptVersion        = errors.New("failed delete prompt version")
	ErrFailedDeletePrompt               = errors.New("failed delete prompt")
	ErrPromptHandleTaken                = errors.New("prompt handle taken")
	ErrInvalidPromptHandle              = errors.New("invalid prompt handle")
	ErrInvalidPromptVersion             = errors.New("invalid prompt version")

	PromptLabelLatest  = "latest"
	promptHandleMaxLen = 100
)

// Prompt is the module for workspace prompt CRUD.
type Prompt struct {
	PromptRepository        db.PromptRepository
	PromptVersionRepository db.PromptVersionRepository
	WorkspaceRepository     db.WorkspaceRepository
}

// CreatePromptRequest is the body for creating a prompt with its first version.
type CreatePromptRequest struct {
	Name          string  `json:"name" validate:"required,min=1,max=150" label:"Name"`
	Description   string  `json:"description" validate:"omitempty,max=2000" label:"Description"`
	Type          string  `json:"type" validate:"required,oneof=text chat" label:"Type"`
	Content       string  `json:"content" validate:"required" label:"Content"`
	Config        *string `json:"config" validate:"omitempty,json" label:"Config"`
	CommitMessage *string `json:"commitMessage" validate:"omitempty,max=255" label:"Commit message"`
	Meta          *string `json:"meta" validate:"omitempty,json" label:"Meta"`
	Production    bool    `json:"production" label:"Production"`
}

// CreatePromptVersionRequest is the body for creating a version on an existing prompt.
type CreatePromptVersionRequest struct {
	Type          string  `json:"type" validate:"required,oneof=text chat" label:"Type"`
	Content       string  `json:"content" validate:"required" label:"Content"`
	Config        *string `json:"config" validate:"omitempty,json" label:"Config"`
	CommitMessage *string `json:"commitMessage" validate:"omitempty,max=255" label:"Commit message"`
	Meta          *string `json:"meta" validate:"omitempty,json" label:"Meta"`
	Production    bool    `json:"production" label:"Production"`
}

// UpdatePromptVersionRequest is the body for updating a prompt version.
type UpdatePromptVersionRequest struct {
	Type          string  `json:"type" validate:"omitempty,oneof=text chat" label:"Type"`
	Content       string  `json:"content" validate:"required" label:"Content"`
	Config        *string `json:"config" validate:"omitempty,json" label:"Config"`
	CommitMessage *string `json:"commitMessage" validate:"omitempty,max=255" label:"Commit message"`
	Meta          *string `json:"meta" validate:"omitempty,json" label:"Meta"`
	Status        string  `json:"status" validate:"omitempty,oneof=active archived" label:"Status"`
	Production    bool    `json:"production" label:"Production"`
}

// PromptVersionResponse is a prompt version shaped for API responses.
type PromptVersionResponse struct {
	Id            db.Id           `json:"id"`
	Name          string          `json:"name"`
	Handle        string          `json:"handle"`
	Description   string          `json:"description"`
	Version       int             `json:"version"`
	Type          string          `json:"type"`
	Content       string          `json:"content"`
	Config        json.RawMessage `json:"config,omitempty"`
	Labels        []string        `json:"labels,omitempty"`
	CommitMessage *string         `json:"commitMessage,omitempty"`
	CommitHash    string          `json:"commitHash"`
	Meta          *string         `json:"meta,omitempty"`
	Status        string          `json:"status"`
	Production    bool            `json:"production"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
}

// ListPromptsResponse is returned when listing prompts.
type ListPromptsResponse struct {
	Prompts []*PromptResponse
	Total   int64
}

// PromptResponse is a prompt for list/get views.
type PromptResponse struct {
	PromptId     string `json:"promptId"`
	Name         string `json:"name"`
	Handle       string `json:"handle"`
	Description  string `json:"description"`
	VersionCount int64  `json:"versionCount"`
}

// ListPromptVersionsResponse is returned when listing prompt versions.
type ListPromptVersionsResponse struct {
	Versions []*PromptVersionResponse
	Total    int64
}

// NewPrompt creates a prompt module with the given repositories.
func NewPrompt(prompts db.PromptRepository, promptVersions db.PromptVersionRepository, workspaces db.WorkspaceRepository) *Prompt {
	return &Prompt{
		PromptRepository:        prompts,
		PromptVersionRepository: promptVersions,
		WorkspaceRepository:     workspaces,
	}
}

// CreatePrompt creates a new prompt with its first version.
func (p *Prompt) CreatePrompt(workspaceId db.Id, req *CreatePromptRequest) (*PromptResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	handle := util.HandleFromName(req.Name, promptHandleMaxLen)
	if lo.IsEmpty(handle) {
		return nil, ErrInvalidPromptHandle
	}
	exists, err := p.PromptRepository.ExistsByHandle(workspaceId, handle)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}
	if exists {
		return nil, ErrPromptHandleTaken
	}

	promptId, err := p.PromptRepository.Create(workspaceId, req.Name, handle, req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}

	_, err = p.CreatePromptVersion(workspaceId, promptId, &CreatePromptVersionRequest{
		Type:          req.Type,
		Content:       req.Content,
		Config:        req.Config,
		CommitMessage: req.CommitMessage,
		Meta:          req.Meta,
		Production:    req.Production,
	})
	if err != nil {
		_ = p.PromptRepository.Delete(workspaceId, promptId)
		return nil, err
	}

	return &PromptResponse{
		PromptId:     promptId.String(),
		Name:         req.Name,
		Handle:       handle,
		Description:  req.Description,
		VersionCount: 1,
	}, nil
}

// CreatePromptVersion creates a new version for an existing prompt.
func (p *Prompt) CreatePromptVersion(workspaceId, promptId db.Id, req *CreatePromptVersionRequest) (*PromptVersionResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	prompt, err := p.PromptRepository.GetById(workspaceId, promptId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}
	if prompt == nil {
		return nil, ErrPromptNotFound
	}

	err = p.PromptVersionRepository.ClearLabelFromPromptId(promptId, PromptLabelLatest)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}

	latestLabels, err := json.Marshal([]string{PromptLabelLatest})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}
	labels := string(latestLabels)

	version := &db.PromptVersion{
		PromptId:      promptId,
		Type:          req.Type,
		Content:       req.Content,
		Config:        req.Config,
		Labels:        &labels,
		CommitMessage: req.CommitMessage,
		Meta:          req.Meta,
		Status:        "active",
		Production:    req.Production,
	}

	err = p.PromptVersionRepository.Create(version)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}

	return &PromptVersionResponse{
		Id:            version.Id,
		Name:          prompt.Name,
		Handle:        prompt.Handle,
		Description:   prompt.Description,
		Version:       version.Version,
		Type:          version.Type,
		Content:       version.Content,
		Config:        util.JSONRawFromString(version.Config),
		Labels:        util.JSONSliceFromString[string](version.Labels),
		CommitMessage: version.CommitMessage,
		CommitHash:    version.CommitHash,
		Meta:          version.Meta,
		Status:        version.Status,
		Production:    version.Production,
		CreatedAt:     version.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     version.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetPromptVersionByHandle returns a prompt version resolved by handle and version selector.
func (p *Prompt) GetPromptVersionByHandle(workspaceId db.Id, handle, versionSpec string) (*PromptVersionResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	if lo.IsEmpty(handle) {
		return nil, ErrInvalidPromptHandle
	}

	prompt, err := p.PromptRepository.GetByHandle(workspaceId, handle)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetPromptVersion, err)
	}
	if prompt == nil {
		return nil, ErrPromptVersionNotFound
	}

	versionSpec = strings.TrimSpace(strings.ToLower(versionSpec))

	var version *db.PromptVersion
	switch versionSpec {
	case "", PromptLabelLatest:
		version, err = p.PromptVersionRepository.GetLatestByPromptId(prompt.Id)
	case "production":
		version, err = p.PromptVersionRepository.GetProductionByPromptId(prompt.Id)
	default:
		vn, parseErr := strconv.Atoi(strings.TrimPrefix(versionSpec, "v"))
		if parseErr != nil || vn < 1 {
			return nil, ErrInvalidPromptVersion
		}
		version, err = p.PromptVersionRepository.GetByPromptIdAndVersion(prompt.Id, vn)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetPromptVersion, err)
	}
	if version == nil {
		return nil, ErrPromptVersionNotFound
	}

	return &PromptVersionResponse{
		Id:            version.Id,
		Name:          prompt.Name,
		Handle:        prompt.Handle,
		Description:   prompt.Description,
		Version:       version.Version,
		Type:          version.Type,
		Content:       version.Content,
		Config:        util.JSONRawFromString(version.Config),
		Labels:        util.JSONSliceFromString[string](version.Labels),
		CommitMessage: version.CommitMessage,
		CommitHash:    version.CommitHash,
		Meta:          version.Meta,
		Status:        version.Status,
		Production:    version.Production,
		CreatedAt:     version.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     version.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetPromptVersion returns one prompt version by Id within a prompt.
func (p *Prompt) GetPromptVersion(workspaceId, promptId, versionId db.Id) (*PromptVersionResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	prompt, err := p.PromptRepository.GetById(workspaceId, promptId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetPromptVersion, err)
	}
	if prompt == nil {
		return nil, ErrPromptVersionNotFound
	}

	version, err := p.PromptVersionRepository.GetById(versionId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetPromptVersion, err)
	}
	if version == nil || version.PromptId != promptId {
		return nil, ErrPromptVersionNotFound
	}

	return &PromptVersionResponse{
		Id:            version.Id,
		Name:          prompt.Name,
		Handle:        prompt.Handle,
		Description:   prompt.Description,
		Version:       version.Version,
		Type:          version.Type,
		Content:       version.Content,
		Config:        util.JSONRawFromString(version.Config),
		Labels:        util.JSONSliceFromString[string](version.Labels),
		CommitMessage: version.CommitMessage,
		CommitHash:    version.CommitHash,
		Meta:          version.Meta,
		Status:        version.Status,
		Production:    version.Production,
		CreatedAt:     version.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     version.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// UpdatePromptVersion updates a prompt version.
func (p *Prompt) UpdatePromptVersion(workspaceId, promptId, versionId db.Id, req *UpdatePromptVersionRequest) (*PromptVersionResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	prompt, err := p.PromptRepository.GetById(workspaceId, promptId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdatePromptVersion, err)
	}
	if prompt == nil {
		return nil, ErrPromptVersionNotFound
	}

	version, err := p.PromptVersionRepository.GetById(versionId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdatePromptVersion, err)
	}
	if version == nil || version.PromptId != promptId {
		return nil, ErrPromptVersionNotFound
	}

	if lo.IsNotEmpty(req.Type) {
		version.Type = req.Type
	}
	version.Content = req.Content
	version.Config = req.Config
	version.CommitMessage = req.CommitMessage
	version.Meta = req.Meta
	if lo.IsNotEmpty(req.Status) {
		version.Status = req.Status
	}
	if req.Production {
		err = p.PromptVersionRepository.ClearProductionFromPromptId(version.PromptId, &versionId)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedSetPromptVersionProduction, err)
		}
		version.Production = true
	}

	err = p.PromptVersionRepository.Update(version)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdatePromptVersion, err)
	}

	return &PromptVersionResponse{
		Id:            version.Id,
		Name:          prompt.Name,
		Handle:        prompt.Handle,
		Description:   prompt.Description,
		Version:       version.Version,
		Type:          version.Type,
		Content:       version.Content,
		Config:        util.JSONRawFromString(version.Config),
		Labels:        util.JSONSliceFromString[string](version.Labels),
		CommitMessage: version.CommitMessage,
		CommitHash:    version.CommitHash,
		Meta:          version.Meta,
		Status:        version.Status,
		Production:    version.Production,
		CreatedAt:     version.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     version.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// DeletePromptVersion deletes a prompt version by Id.
func (p *Prompt) DeletePromptVersion(workspaceId, promptId, versionId db.Id) error {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	prompt, err := p.PromptRepository.GetById(workspaceId, promptId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeletePromptVersion, err)
	}
	if prompt == nil {
		return ErrPromptVersionNotFound
	}

	version, err := p.PromptVersionRepository.GetById(versionId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeletePromptVersion, err)
	}
	if version == nil || version.PromptId != promptId {
		return ErrPromptVersionNotFound
	}

	err = p.PromptVersionRepository.Delete(versionId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeletePromptVersion, err)
	}

	return nil
}

// ListPrompts returns paginated prompt summaries for a workspace.
func (p *Prompt) ListPrompts(workspaceId db.Id, limit, offset int) (*ListPromptsResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	total, err := p.PromptRepository.CountByWorkspaceId(workspaceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPrompts, err)
	}

	results, err := p.PromptRepository.ListByWorkspace(workspaceId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPrompts, err)
	}

	list := make([]*PromptResponse, 0, len(results))
	for _, item := range results {
		list = append(list, &PromptResponse{
			PromptId:     item.Id.String(),
			Name:         item.Name,
			Handle:       item.Handle,
			Description:  item.Description,
			VersionCount: item.VersionCount,
		})
	}

	return &ListPromptsResponse{Prompts: list, Total: total}, nil
}

// GetPrompt returns one prompt by Id.
func (p *Prompt) GetPrompt(workspaceId, promptId db.Id) (*PromptResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	result, err := p.PromptRepository.GetById(workspaceId, promptId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetPrompt, err)
	}
	if result == nil {
		return nil, ErrPromptNotFound
	}

	return &PromptResponse{
		PromptId:     result.Id.String(),
		Name:         result.Name,
		Handle:       result.Handle,
		Description:  result.Description,
		VersionCount: result.VersionCount,
	}, nil
}

// ListPromptVersions returns paginated versions for a prompt.
func (p *Prompt) ListPromptVersions(workspaceId, promptId db.Id, limit, offset int) (*ListPromptVersionsResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	prompt, err := p.PromptRepository.GetById(workspaceId, promptId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPromptVersions, err)
	}
	if prompt == nil {
		return nil, ErrPromptNotFound
	}

	total, err := p.PromptVersionRepository.CountByPromptId(promptId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPromptVersions, err)
	}

	versions, err := p.PromptVersionRepository.ListByPromptId(promptId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPromptVersions, err)
	}

	list := make([]*PromptVersionResponse, 0, len(versions))
	for _, version := range versions {
		list = append(list, &PromptVersionResponse{
			Id:            version.Id,
			Name:          prompt.Name,
			Handle:        prompt.Handle,
			Description:   prompt.Description,
			Version:       version.Version,
			Type:          version.Type,
			Content:       version.Content,
			Config:        util.JSONRawFromString(version.Config),
			Labels:        util.JSONSliceFromString[string](version.Labels),
			CommitMessage: version.CommitMessage,
			CommitHash:    version.CommitHash,
			Meta:          version.Meta,
			Status:        version.Status,
			Production:    version.Production,
			CreatedAt:     version.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:     version.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return &ListPromptVersionsResponse{Versions: list, Total: total}, nil
}

// DeletePrompt deletes a prompt and all of its versions.
func (p *Prompt) DeletePrompt(workspaceId, promptId db.Id) error {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	err = p.PromptRepository.Delete(workspaceId, promptId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPromptNotFound
		}
		return fmt.Errorf("%w: %v", ErrFailedDeletePrompt, err)
	}

	return nil
}
