// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/clivern/ziee/conf"
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/locale"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/ai"
	"github.com/clivern/ziee/pkg/qdrant"
	"github.com/clivern/ziee/pkg/util"
	"github.com/clivern/ziee/service/knowledge"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// UploadDocumentAction uploads a .txt or .md document to a workspace.
func (a *API) UploadDocumentAction(w http.ResponseWriter, r *http.Request) {
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

	doc, err := a.Document.UploadDocument(
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
func (a *API) ListDocumentsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	limit, offset := util.ParsePagination(r)

	result, err := a.Document.ListDocuments(
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
func (a *API) DeleteDocumentAction(w http.ResponseWriter, r *http.Request) {
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

	err := a.Document.DeleteDocument(
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
func (a *API) SearchDocumentsAction(w http.ResponseWriter, r *http.Request) {
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

	dm := &module.Document{
		DocumentRepository:     a.Document.DocumentRepository,
		WorkspaceRepository:    a.Document.WorkspaceRepository,
		UsageRepository:        a.Usage,
		SubscriptionRepository: a.Subscriptions,
		Knowledge: knowledge.New(knowledge.Dependencies{
			Documents:     a.Documents,
			Embed:         ai.NewEmbedClient(),
			Vectors:       vdb,
			Store:         a.Store,
			Usage:         a.Usage,
			Subscriptions: a.Subscriptions,
		}),
	}

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
