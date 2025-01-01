// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/actx0/ziee/conf"
	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/ai"
	"github.com/actx0/ziee/pkg/qdrant"
	"github.com/actx0/ziee/pkg/storage"
	"github.com/actx0/ziee/pkg/util"
	"github.com/actx0/ziee/service/knowledge"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// UploadDocumentAction uploads a .txt or .md document to a workspace.
func UploadDocumentAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	form, err := util.ParseUploadForm(r)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": err.Error(),
		})
		return
	}

	store, err := storage.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize document storage")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_process_request"),
		})
		return
	}

	dm := module.NewDocument(
		db.NewWorkspaceDocumentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		store,
	)

	doc, err := dm.UploadDocument(
		r.Context(),
		form,
		db.Id(wid),
	)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Msg("Failed to upload document")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_upload_document"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, doc)
}

// ListDocumentsAction returns documents for a workspace.
func ListDocumentsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	limit, offset := util.ParsePagination(r)

	store, err := storage.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize document storage")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_process_request"),
		})
		return
	}

	dm := module.NewDocument(
		db.NewWorkspaceDocumentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		store,
	)

	result, err := dm.ListDocuments(
		db.Id(wid),
		limit,
		offset,
	)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Msg("Failed to list documents")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_documents"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"documents": result.Documents,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// DeleteDocumentAction deletes a workspace document.
func DeleteDocumentAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	documentId := chi.URLParam(r, "documentId")
	if lo.IsEmpty(documentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_document_id"),
		})
		return
	}

	store, err := storage.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize document storage")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_process_request"),
		})
		return
	}

	dm := module.NewDocument(
		db.NewWorkspaceDocumentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		store,
	)

	err = dm.DeleteDocument(
		r.Context(),
		db.Id(wid),
		db.Id(documentId),
	)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrDocumentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "document_not_found"),
			})
		default:
			log.Error().Err(err).Str("documentId", documentId).Msg("Failed to delete document")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_document"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchDocumentsAction searches workspace documents by semantic query.
func SearchDocumentsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	var req module.SearchDocumentsRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	limit := lo.Ternary(req.Limit == 0, conf.DefaultSearchLimit, req.Limit)

	store, err := storage.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize document storage")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_process_request"),
		})
		return
	}

	vdb, err := qdrant.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize qdrant")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_process_request"),
		})
		return
	}
	defer func() {
		if err := vdb.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close qdrant client")
		}
	}()

	ksvc := knowledge.New(knowledge.Dependencies{
		Documents:     db.NewWorkspaceDocumentRepository(db.GetDB()),
		Embed:         ai.NewEmbedClient(),
		Vectors:       vdb,
		Store:         store,
		Usage:         db.NewUsageRepository(db.GetDB()),
		Subscriptions: db.NewSubscriptionRepository(db.GetDB()),
	})

	dm := module.NewDocument(
		db.NewWorkspaceDocumentRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		nil,
	)
	dm.Knowledge = ksvc
	dm.UsageRepository = db.NewUsageRepository(db.GetDB())
	dm.SubscriptionRepository = db.NewSubscriptionRepository(db.GetDB())

	result, err := dm.SearchDocuments(
		r.Context(),
		db.Id(wid),
		req.Query,
		req.Labels,
		limit,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Msg("Failed to search documents")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_search_documents"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, result)
}
