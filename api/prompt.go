// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// CreatePromptAction creates a new prompt with its first version.
func CreatePromptAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	var req module.CreatePromptRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	pm := module.NewPrompt(
		db.NewPromptRepository(db.GetDB()),
		db.NewPromptVersionRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	version, err := pm.CreatePrompt(workspaceId, &req)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrPromptHandleTaken):
			util.WriteJSON(w, http.StatusConflict, map[string]any{
				"errorMessage": locale.TR(r, "prompt_handle_taken"),
			})
		case errors.Is(err, module.ErrInvalidPromptHandle):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_prompt_handle"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", workspaceId.String()).Msg("Failed to create prompt")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_prompt_version"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, version)
}

// CreatePromptVersionAction creates a new version for an existing prompt.
func CreatePromptVersionAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	promptId := chi.URLParam(r, "promptId")
	if lo.IsEmpty(promptId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_id"),
		})
		return
	}

	var req module.CreatePromptVersionRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	pm := module.NewPrompt(
		db.NewPromptRepository(db.GetDB()),
		db.NewPromptVersionRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	version, err := pm.CreatePromptVersion(workspaceId, db.Id(promptId), &req)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrPromptNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "prompt_not_found"),
			})
		default:
			log.Error().Err(err).Str("promptId", promptId).Msg("Failed to create prompt version")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_prompt_version"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, version)
}

// ListPromptsAction returns prompt summaries for a workspace (one row per prompt name).
func ListPromptsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	limit, offset := util.ParsePagination(r)

	pm := module.NewPrompt(
		db.NewPromptRepository(db.GetDB()),
		db.NewPromptVersionRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	result, err := pm.ListPrompts(workspaceId, limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", workspaceId.String()).Msg("Failed to list prompts")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_prompts"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"prompts": result.Prompts,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetPromptByNameAction returns a prompt version by handle and optional version selector.
func GetPromptByNameAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	name := chi.URLParam(r, "promptName")
	if lo.IsEmpty(name) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_name"),
		})
		return
	}

	pm := module.NewPrompt(
		db.NewPromptRepository(db.GetDB()),
		db.NewPromptVersionRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	promptVersion, err := pm.GetPromptVersionByHandle(workspaceId, name, r.URL.Query().Get("version"))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrInvalidPromptHandle):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_prompt_name"),
			})
		case errors.Is(err, module.ErrInvalidPromptVersion):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_prompt_version"),
			})
		case errors.Is(err, module.ErrPromptVersionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "prompt_version_not_found"),
			})
		default:
			log.Error().Err(err).Str("promptName", name).Msg("Failed to get prompt version by name")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_prompt_version"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, promptVersion)
}

// GetPromptAction returns one prompt by Id.
func GetPromptAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	promptId := chi.URLParam(r, "promptId")
	if lo.IsEmpty(promptId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_id"),
		})
		return
	}

	pm := module.NewPrompt(
		db.NewPromptRepository(db.GetDB()),
		db.NewPromptVersionRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	prompt, err := pm.GetPrompt(workspaceId, db.Id(promptId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrPromptNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "prompt_not_found"),
			})
		default:
			log.Error().Err(err).Str("promptId", promptId).Msg("Failed to get prompt")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_prompt"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, prompt)
}

// GetPromptVersionAction returns one prompt version by Id.
func GetPromptVersionAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	promptId := chi.URLParam(r, "promptId")
	if lo.IsEmpty(promptId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_id"),
		})
		return
	}

	promptVersionId := chi.URLParam(r, "promptVersionId")
	if lo.IsEmpty(promptVersionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_version_id"),
		})
		return
	}

	pm := module.NewPrompt(
		db.NewPromptRepository(db.GetDB()),
		db.NewPromptVersionRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	version, err := pm.GetPromptVersion(workspaceId, db.Id(promptId), db.Id(promptVersionId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrPromptVersionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "prompt_version_not_found"),
			})
		default:
			log.Error().Err(err).Str("versionId", promptVersionId).Msg("Failed to get prompt version")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_prompt_version"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, version)
}

// UpdatePromptVersionAction updates a prompt version.
func UpdatePromptVersionAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	promptId := chi.URLParam(r, "promptId")
	if lo.IsEmpty(promptId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_id"),
		})
		return
	}

	promptVersionId := chi.URLParam(r, "promptVersionId")
	if lo.IsEmpty(promptVersionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_version_id"),
		})
		return
	}

	var req module.UpdatePromptVersionRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	pm := module.NewPrompt(
		db.NewPromptRepository(db.GetDB()),
		db.NewPromptVersionRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	version, err := pm.UpdatePromptVersion(
		workspaceId,
		db.Id(promptId),
		db.Id(promptVersionId),
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrPromptVersionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "prompt_version_not_found"),
			})
		default:
			log.Error().Err(err).Str("versionId", promptVersionId).Msg("Failed to update prompt version")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_update_prompt_version"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, version)
}

// DeletePromptVersionAction deletes a prompt version by Id.
func DeletePromptVersionAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	promptId := chi.URLParam(r, "promptId")
	if lo.IsEmpty(promptId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_id"),
		})
		return
	}

	promptVersionId := chi.URLParam(r, "promptVersionId")
	if lo.IsEmpty(promptVersionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_version_id"),
		})
		return
	}

	pm := module.NewPrompt(
		db.NewPromptRepository(db.GetDB()),
		db.NewPromptVersionRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	err := pm.DeletePromptVersion(workspaceId, db.Id(promptId), db.Id(promptVersionId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrPromptVersionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "prompt_version_not_found"),
			})
		default:
			log.Error().Err(err).Str("versionId", promptVersionId).Msg("Failed to delete prompt version")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_prompt_version"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListPromptVersionsAction returns versions for a prompt.
func ListPromptVersionsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	promptId := chi.URLParam(r, "promptId")
	if lo.IsEmpty(promptId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_id"),
		})
		return
	}

	limit, offset := util.ParsePagination(r)

	pm := module.NewPrompt(
		db.NewPromptRepository(db.GetDB()),
		db.NewPromptVersionRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	result, err := pm.ListPromptVersions(workspaceId, db.Id(promptId), limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrPromptNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "prompt_not_found"),
			})
		default:
			log.Error().Err(err).Str("promptId", promptId).Msg("Failed to list prompt versions")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_prompts"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"versions": result.Versions,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// DeletePromptAction deletes a prompt and all of its versions.
func DeletePromptAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	promptId := chi.URLParam(r, "promptId")
	if lo.IsEmpty(promptId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_prompt_id"),
		})
		return
	}

	pm := module.NewPrompt(
		db.NewPromptRepository(db.GetDB()),
		db.NewPromptVersionRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	err := pm.DeletePrompt(workspaceId, db.Id(promptId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrPromptNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "prompt_not_found"),
			})
		default:
			log.Error().Err(err).Str("promptId", promptId).Msg("Failed to delete prompt")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_prompt"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
