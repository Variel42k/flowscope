package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/flowscope/flowscope/internal/auth"
	"github.com/flowscope/flowscope/internal/model"
)

func (s *Server) handleAlertRulesList(w http.ResponseWriter, r *http.Request) {
	items, err := s.store.ListAlertRules(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (s *Server) handleAlertRulesCreate(w http.ResponseWriter, r *http.Request) {
	var req model.AlertRule
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.RuleType = strings.TrimSpace(req.RuleType)
	if req.Name == "" || req.RuleType == "" {
		writeErr(w, http.StatusBadRequest, errBadRequest("name and rule_type are required"))
		return
	}
	user := auth.AuthUserFromContext(r.Context())
	req.CreatedBy = user.Username
	item, err := s.store.CreateAlertRule(r.Context(), req)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) handleAlertRulesUpdate(w http.ResponseWriter, r *http.Request) {
	ruleID := chi.URLParam(r, "id")
	if strings.TrimSpace(ruleID) == "" {
		writeErr(w, http.StatusBadRequest, errBadRequest("missing rule id"))
		return
	}
	var req model.AlertRule
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	req.RuleID = ruleID
	item, err := s.store.UpdateAlertRule(r.Context(), ruleID, req)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleAlertRulesDelete(w http.ResponseWriter, r *http.Request) {
	ruleID := chi.URLParam(r, "id")
	if strings.TrimSpace(ruleID) == "" {
		writeErr(w, http.StatusBadRequest, errBadRequest("missing rule id"))
		return
	}
	user := auth.AuthUserFromContext(r.Context())
	if err := s.store.DeleteAlertRule(r.Context(), ruleID, user.Username); err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": ruleID})
}

func (s *Server) handleAlertEventsList(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 50
	}
	if f.PageSize > 200 {
		f.PageSize = 200
	}
	if strings.TrimSpace(r.URL.Query().Get("from")) == "" {
		f.From = time.Now().UTC().Add(-24 * time.Hour)
	}
	if strings.TrimSpace(r.URL.Query().Get("to")) == "" {
		f.To = time.Now().UTC()
	}
	severity := strings.TrimSpace(r.URL.Query().Get("severity"))
	limit := parseInt(r.URL.Query().Get("limit"), 0)
	if limit > 0 {
		f.Page = 1
		f.PageSize = limit
	}
	items, err := s.store.ListAlertEvents(r.Context(), f, severity)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleAlertsEvaluateNow(w http.ResponseWriter, r *http.Request) {
	refTs := time.Now().UTC()
	if v := strings.TrimSpace(r.URL.Query().Get("at")); v != "" {
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			refTs = parsed.UTC()
		}
	}
	events, err := s.store.EvaluateAlerts(r.Context(), refTs)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	inserted, err := s.store.StoreAlertEvents(r.Context(), events)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"generated": len(events),
		"inserted":  inserted,
		"at":        strconv.FormatInt(refTs.Unix(), 10),
	})
}
