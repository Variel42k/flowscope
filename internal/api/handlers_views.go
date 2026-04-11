package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/flowscope/flowscope/internal/auth"
	"github.com/flowscope/flowscope/internal/model"
)

func (s *Server) handleViewsList(w http.ResponseWriter, r *http.Request) {
	scope := strings.TrimSpace(r.URL.Query().Get("scope"))
	user := auth.AuthUserFromContext(r.Context())
	items, err := s.store.ListSavedViews(r.Context(), user.Username, scope)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (s *Server) handleViewsCreate(w http.ResponseWriter, r *http.Request) {
	var req model.SavedView
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	user := auth.AuthUserFromContext(r.Context())
	if strings.TrimSpace(req.Name) == "" {
		writeErr(w, http.StatusBadRequest, errBadRequest("name is required"))
		return
	}
	if req.Filters == nil {
		req.Filters = map[string]string{}
	}
	req.Scope = normalizeScope(req.Scope)
	req.OwnerUser = user.Username
	if req.IsShared && !strings.EqualFold(user.Role, s.cfg.AdminRole) {
		writeErr(w, http.StatusForbidden, errBadRequest("only admins can create shared views"))
		return
	}
	req.CreatedAt = time.Now().UTC()
	req.UpdatedAt = req.CreatedAt
	item, err := s.store.CreateSavedView(r.Context(), req)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) handleViewsUpdate(w http.ResponseWriter, r *http.Request) {
	viewID := strings.TrimSpace(chi.URLParam(r, "id"))
	if viewID == "" {
		writeErr(w, http.StatusBadRequest, errBadRequest("missing view id"))
		return
	}
	existing, err := s.store.GetSavedViewByID(r.Context(), viewID)
	if err != nil {
		writeErr(w, http.StatusNotFound, err)
		return
	}
	user := auth.AuthUserFromContext(r.Context())
	isAdmin := strings.EqualFold(user.Role, s.cfg.AdminRole)
	if !isAdmin && !strings.EqualFold(existing.OwnerUser, user.Username) {
		writeErr(w, http.StatusForbidden, errBadRequest("cannot modify another user's view"))
		return
	}
	var req model.SavedView
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if strings.TrimSpace(req.Name) != "" {
		existing.Name = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		existing.Description = strings.TrimSpace(req.Description)
	}
	if strings.TrimSpace(req.Scope) != "" {
		existing.Scope = normalizeScope(req.Scope)
	}
	if req.Filters != nil {
		existing.Filters = req.Filters
	}
	if req.IsShared != existing.IsShared {
		if !isAdmin {
			writeErr(w, http.StatusForbidden, errBadRequest("only admins can change shared flag"))
			return
		}
		existing.IsShared = req.IsShared
	}
	existing.UpdatedAt = time.Now().UTC()
	item, err := s.store.UpdateSavedView(r.Context(), existing)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleViewsDelete(w http.ResponseWriter, r *http.Request) {
	viewID := strings.TrimSpace(chi.URLParam(r, "id"))
	if viewID == "" {
		writeErr(w, http.StatusBadRequest, errBadRequest("missing view id"))
		return
	}
	existing, err := s.store.GetSavedViewByID(r.Context(), viewID)
	if err != nil {
		writeErr(w, http.StatusNotFound, err)
		return
	}
	user := auth.AuthUserFromContext(r.Context())
	isAdmin := strings.EqualFold(user.Role, s.cfg.AdminRole)
	if !isAdmin && !strings.EqualFold(existing.OwnerUser, user.Username) {
		writeErr(w, http.StatusForbidden, errBadRequest("cannot delete another user's view"))
		return
	}
	if err := s.store.DeleteSavedView(r.Context(), viewID); err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": viewID})
}

func normalizeScope(scope string) string {
	scope = strings.ToLower(strings.TrimSpace(scope))
	switch scope {
	case "overview", "flows", "sankey", "map", "global":
		return scope
	default:
		return "global"
	}
}
