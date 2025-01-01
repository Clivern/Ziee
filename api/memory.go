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

// CreateSessionMemoryAction creates a memory in an agent session.
func CreateSessionMemoryAction(w http.ResponseWriter, r *http.Request) {
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

	var req module.CreateSessionMemoryRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	memory, err := mm.CreateSessionMemory(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
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
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to create session memory")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_session_memory"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, memory)
}

// BatchCreateSessionMemoriesAction creates multiple memories in an agent session.
func BatchCreateSessionMemoriesAction(w http.ResponseWriter, r *http.Request) {
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

	var req module.BatchCreateSessionMemoriesRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	memories, err := mm.BatchCreateSessionMemories(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
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
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to batch create session memories")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_session_memory"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, map[string]any{
		"memories": memories,
	})
}

// ListSessionMemoriesAction lists memories for an agent session.
func ListSessionMemoriesAction(w http.ResponseWriter, r *http.Request) {
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

	limit, offset := util.ParsePagination(r)

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	result, err := mm.ListSessionMemories(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
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
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to list session memories")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_session_memories"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"memories": result.Memories,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetSessionMemoryAction returns one memory by id.
func GetSessionMemoryAction(w http.ResponseWriter, r *http.Request) {
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

	memoryId := chi.URLParam(r, "memoryId")
	if lo.IsEmpty(memoryId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_memory_id"),
		})
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	memory, err := mm.GetSessionMemory(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		db.Id(memoryId),
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
		case errors.Is(err, module.ErrSessionMemoryNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "session_memory_not_found"),
			})
		default:
			log.Error().Err(err).Str("memoryId", memoryId).Msg("Failed to get session memory")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_session_memory"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, memory)
}

// UpdateSessionMemoryAction updates a memory by id.
func UpdateSessionMemoryAction(w http.ResponseWriter, r *http.Request) {
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

	memoryId := chi.URLParam(r, "memoryId")
	if lo.IsEmpty(memoryId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_memory_id"),
		})
		return
	}

	var req module.UpdateSessionMemoryRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	memory, err := mm.UpdateSessionMemory(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		db.Id(memoryId),
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
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrSessionMemoryNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "session_memory_not_found"),
			})
		default:
			log.Error().Err(err).Str("memoryId", memoryId).Msg("Failed to update session memory")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_update_session_memory"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, memory)
}

// DeleteSessionMemoryAction deletes a memory by id.
func DeleteSessionMemoryAction(w http.ResponseWriter, r *http.Request) {
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

	memoryId := chi.URLParam(r, "memoryId")
	if lo.IsEmpty(memoryId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_memory_id"),
		})
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	err := mm.DeleteSessionMemory(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		db.Id(memoryId),
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
		case errors.Is(err, module.ErrSessionMemoryNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "session_memory_not_found"),
			})
		default:
			log.Error().Err(err).Str("memoryId", memoryId).Msg("Failed to delete session memory")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_session_memory"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// BatchDeleteSessionMemoriesAction deletes multiple memories in an agent session.
func BatchDeleteSessionMemoriesAction(w http.ResponseWriter, r *http.Request) {
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

	var req module.BatchDeleteSessionMemoriesRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	err := mm.BatchDeleteSessionMemories(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
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
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrSessionMemoryNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "session_memory_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to batch delete session memories")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_session_memory"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
