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

// CreateAgentSessionAction creates an agent session identified by query id and/or labels.
func CreateAgentSessionAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	externalId := r.URL.Query().Get("id")
	labels := util.ParseQueryLabels(r)
	if lo.IsEmpty(externalId) && len(labels) == 0 {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	var req module.CreateAgentSessionRequest
	if r.ContentLength > 0 {
		if err := util.DecodeAndValidate(r, &req); err != nil {
			util.WriteValidationError(w, err)
			return
		}
	}

	sm := module.NewAgentSession(
		db.NewAgentSessionRepository(db.GetDB()),
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	session, err := sm.CreateAgentSession(
		workspaceId,
		db.Id(agentId),
		externalId,
		labels,
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrInvalidSessionReference):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_session_reference"),
			})
		case errors.Is(err, module.ErrInvalidSessionLabels):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_session_labels"),
			})
		case errors.Is(err, module.ErrAgentSessionAlreadyExists):
			util.WriteJSON(w, http.StatusConflict, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_already_exists"),
			})
		default:
			log.Error().Err(err).Str("agentId", agentId).Msg("Failed to create agent session")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_agent_session"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, session)
}

// ListAgentSessionsAction lists agent sessions with optional id or label filters.
func ListAgentSessionsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	limit, offset := util.ParsePagination(r)
	externalId := r.URL.Query().Get("id")
	labels := util.ParseQueryLabels(r)

	sm := module.NewAgentSession(
		db.NewAgentSessionRepository(db.GetDB()),
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	result, err := sm.ListAgentSessions(
		workspaceId,
		db.Id(agentId),
		externalId,
		labels,
		limit,
		offset,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrInvalidSessionLabels):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_session_labels"),
			})
		default:
			log.Error().Err(err).Str("agentId", agentId).Msg("Failed to list agent sessions")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_agent_sessions"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"sessions": result.Sessions,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetAgentSessionAction returns one agent session by internal id.
func GetAgentSessionAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	sm := module.NewAgentSession(
		db.NewAgentSessionRepository(db.GetDB()),
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	session, err := sm.GetAgentSession(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to get agent session")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_agent_session"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, session)
}

// GetAgentSessionByLabelsAction returns one agent session by external id or exact query labels.
func GetAgentSessionByLabelsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	externalId := r.URL.Query().Get("id")
	labels := util.ParseQueryLabels(r)
	if lo.IsEmpty(externalId) && len(labels) == 0 {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	sm := module.NewAgentSession(
		db.NewAgentSessionRepository(db.GetDB()),
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	session, err := sm.GetAgentSessionByLabels(
		workspaceId,
		db.Id(agentId),
		externalId,
		labels,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrInvalidSessionReference):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_session_reference"),
			})
		case errors.Is(err, module.ErrInvalidSessionLabels):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_session_labels"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("agentId", agentId).Msg("Failed to get agent session")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_agent_session"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, session)
}

// UpdateAgentSessionByLabelsAction updates an agent session identified by external id or query labels.
func UpdateAgentSessionByLabelsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	externalId := r.URL.Query().Get("id")
	labels := util.ParseQueryLabels(r)
	if lo.IsEmpty(externalId) && len(labels) == 0 {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	var req module.UpdateAgentSessionRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	sm := module.NewAgentSession(
		db.NewAgentSessionRepository(db.GetDB()),
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	session, err := sm.UpdateAgentSessionByLabels(
		workspaceId,
		db.Id(agentId),
		externalId,
		labels,
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrInvalidSessionReference):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_session_reference"),
			})
		case errors.Is(err, module.ErrInvalidSessionLabels):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_session_labels"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionAlreadyExists):
			util.WriteJSON(w, http.StatusConflict, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_already_exists"),
			})
		default:
			log.Error().Err(err).Str("agentId", agentId).Msg("Failed to update agent session")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_update_agent_session"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, session)
}

// DeleteAgentSessionByLabelsAction deletes an agent session identified by external id or query labels.
func DeleteAgentSessionByLabelsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	externalId := r.URL.Query().Get("id")
	labels := util.ParseQueryLabels(r)
	if lo.IsEmpty(externalId) && len(labels) == 0 {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	sm := module.NewAgentSession(
		db.NewAgentSessionRepository(db.GetDB()),
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	err := sm.DeleteAgentSessionByLabels(
		workspaceId,
		db.Id(agentId),
		externalId,
		labels,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrInvalidSessionReference):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_session_reference"),
			})
		case errors.Is(err, module.ErrInvalidSessionLabels):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_session_labels"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("agentId", agentId).Msg("Failed to delete agent session")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_agent_session"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
