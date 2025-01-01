// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/locale"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// CreateAgentAction creates a new agent in a workspace.
func CreateAgentAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	var req module.CreateAgentRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	am := module.NewAgent(
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	agent, err := am.CreateAgent(workspaceId, &req)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", workspaceId.String()).Msg("Failed to create agent")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_agent"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, agent)
}

// ListAgentsAction returns agents for a workspace.
func ListAgentsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	limit, offset := util.ParsePagination(r)

	am := module.NewAgent(
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	result, err := am.ListAgents(workspaceId, limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", workspaceId.String()).Msg("Failed to list agents")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_agents"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"agents": result.Agents,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetAgentAction returns one agent by Id.
func GetAgentAction(w http.ResponseWriter, r *http.Request) {
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

	am := module.NewAgent(
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	agent, err := am.GetAgent(workspaceId, db.Id(agentId))
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
		default:
			log.Error().Err(err).Str("agentId", agentId).Msg("Failed to get agent")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_agent"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, agent)
}

// UpdateAgentAction updates an agent by Id.
func UpdateAgentAction(w http.ResponseWriter, r *http.Request) {
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

	var req module.UpdateAgentRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	am := module.NewAgent(
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	agent, err := am.UpdateAgent(workspaceId, db.Id(agentId), &req)
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
		default:
			log.Error().Err(err).Str("agentId", agentId).Msg("Failed to update agent")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_update_agent"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, agent)
}

// DeleteAgentAction deletes an agent by Id.
func DeleteAgentAction(w http.ResponseWriter, r *http.Request) {
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

	am := module.NewAgent(
		db.NewAgentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	err := am.DeleteAgent(workspaceId, db.Id(agentId))
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
		default:
			log.Error().Err(err).Str("agentId", agentId).Msg("Failed to delete agent")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_agent"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
