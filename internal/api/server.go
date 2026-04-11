package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/flowscope/flowscope/internal/auth"
	"github.com/flowscope/flowscope/internal/config"
	graphpkg "github.com/flowscope/flowscope/internal/graph"
	"github.com/flowscope/flowscope/internal/model"
	"github.com/flowscope/flowscope/internal/storage"
)

type Server struct {
	cfg     config.Config
	store   *storage.Repository
	auth    *auth.Manager
	oidc    *auth.OIDCProvider
	handler http.Handler
	started time.Time
}

func NewServer(cfg config.Config, store *storage.Repository, authMgr *auth.Manager) *Server {
	s := &Server{
		cfg:   cfg,
		store: store,
		auth:  authMgr,
		oidc: auth.NewOIDCProvider(auth.OIDCConfig{
			Enabled:            cfg.OIDCEnabled,
			IssuerURL:          cfg.OIDCIssuerURL,
			ClientID:           cfg.OIDCClientID,
			ClientSecret:       cfg.OIDCClientSecret,
			RedirectURL:        cfg.OIDCRedirectURL,
			SuccessRedirectURL: cfg.OIDCSuccessRedirect,
			Scopes:             splitCSV(cfg.OIDCScopes),
			RoleClaim:          cfg.OIDCRoleClaim,
			AdminUsers:         cfg.OIDCAdminUsers,
			AllowedDomains:     cfg.OIDCAllowedDomains,
			AdminRole:          cfg.AdminRole,
			ViewerRole:         cfg.ViewerRole,
		}),
		started: time.Now().UTC(),
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)
	r.Get("/api/health", s.handleHealth)
	r.Post("/api/auth/login", s.handleLogin)
	r.Get("/api/auth/oidc/start", s.handleOIDCStart)
	r.Get("/api/auth/oidc/callback", s.handleOIDCCallback)

	r.Group(func(protected chi.Router) {
		protected.Use(auth.Middleware(authMgr, cfg.AuthDisabled))
		protected.Get("/api/auth/me", s.handleAuthMe)
		protected.Get("/api/exporters", s.handleExporters)
		protected.Get("/api/flows/active", s.handleFlowsActive)
		protected.Get("/api/flows/historical", s.handleFlowsHistorical)
		protected.Get("/api/flows/{id}", s.handleFlowByID)
		protected.Get("/api/talkers/top", s.handleTopTalkers)
		protected.Get("/api/protocols/top", s.handleTopProtocols)
		protected.Get("/api/interfaces/top", s.handleTopInterfaces)
		protected.Get("/api/ports/top", s.handleTopPorts)
		protected.Get("/api/sankey", s.handleSankey)
		protected.Get("/api/map/graph", s.handleGraph)
		protected.Get("/api/map/node/{id}", s.handleMapNode)
		protected.Get("/api/map/edge/{id}", s.handleMapEdge)
		protected.Get("/api/search", s.handleSearch)
		protected.Get("/api/alerts/rules", s.handleAlertRulesList)
		protected.Get("/api/alerts/events", s.handleAlertEventsList)
		protected.Post("/api/alerts/evaluate", s.handleAlertsEvaluateNow)
		protected.Get("/api/views", s.handleViewsList)
		protected.Post("/api/views", s.handleViewsCreate)
		protected.Put("/api/views/{id}", s.handleViewsUpdate)
		protected.Delete("/api/views/{id}", s.handleViewsDelete)
		protected.With(auth.RequireRole(cfg.AdminRole)).Post("/api/inventory/import", s.handleInventoryImport)
		protected.With(auth.RequireRole(cfg.AdminRole)).Post("/api/alerts/rules", s.handleAlertRulesCreate)
		protected.With(auth.RequireRole(cfg.AdminRole)).Put("/api/alerts/rules/{id}", s.handleAlertRulesUpdate)
		protected.With(auth.RequireRole(cfg.AdminRole)).Delete("/api/alerts/rules/{id}", s.handleAlertRulesDelete)
	})
	s.handler = r
	return s
}

func (s *Server) Handler() http.Handler {
	return s.handler
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.store.Health(r.Context()); err != nil {
		writeErr(w, http.StatusServiceUnavailable, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":       "ok",
		"service":      "flowscope-api",
		"uptime_sec":   int(time.Since(s.started).Seconds()),
		"auth_enabled": !s.cfg.AuthDisabled,
		"oidc_enabled": s.cfg.OIDCEnabled,
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if req.Username != s.cfg.AdminUser || req.Password != s.cfg.AdminPassword {
		writeErr(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}
	role := s.cfg.AdminRole
	if strings.TrimSpace(role) == "" {
		role = "admin"
	}
	token, err := s.auth.Sign(req.Username, role)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": req.Username, "role": role, "auth_mode": "local"})
}

func (s *Server) handleExporters(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	data, err := s.store.ListExporters(r.Context(), f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleFlowsActive(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	data, err := s.store.QueryActiveFlows(r.Context(), f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleFlowsHistorical(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	data, err := s.store.QueryHistoricalFlows(r.Context(), f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleFlowByID(w http.ResponseWriter, r *http.Request) {
	flowID := chi.URLParam(r, "id")
	item, err := s.store.GetFlowByID(r.Context(), flowID)
	if err != nil {
		writeErr(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleTopTalkers(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	limit := parseInt(r.URL.Query().Get("limit"), 10)
	data, err := s.store.TopTalkers(r.Context(), f, limit)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": data})
}

func (s *Server) handleTopProtocols(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	limit := parseInt(r.URL.Query().Get("limit"), 10)
	data, err := s.store.TopProtocols(r.Context(), f, limit)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": data})
}

func (s *Server) handleTopInterfaces(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	limit := parseInt(r.URL.Query().Get("limit"), 10)
	data, err := s.store.TopInterfaces(r.Context(), f, limit)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": data})
}

func (s *Server) handleSankey(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	data, err := s.store.Sankey(r.Context(), f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleGraph(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	if f.Mode == "" {
		f.Mode = "host_to_host"
	}
	if f.GroupBy == "" {
		f.GroupBy = "none"
	}
	graph, err := s.store.Graph(r.Context(), f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	graph = graphpkg.Collapse(graph, f.GroupBy)
	if strings.TrimSpace(f.NodeID) != "" {
		graph = isolateEgo(graph, f.NodeID)
	}
	writeJSON(w, http.StatusOK, graph)
}

func isolateEgo(graph model.GraphResponse, nodeID string) model.GraphResponse {
	keepNodes := map[string]bool{nodeID: true}
	keepEdges := make([]model.GraphEdge, 0, len(graph.Edges))
	for _, e := range graph.Edges {
		if e.Source == nodeID || e.Destination == nodeID {
			keepEdges = append(keepEdges, e)
			keepNodes[e.Source] = true
			keepNodes[e.Destination] = true
		}
	}
	nodes := make([]model.GraphNode, 0, len(graph.Nodes))
	for _, n := range graph.Nodes {
		if keepNodes[n.ID] {
			nodes = append(nodes, n)
		}
	}
	graph.Nodes = nodes
	graph.Edges = keepEdges
	return graph
}

func (s *Server) handleMapNode(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "id")
	f := parseFilter(r)
	data, err := s.store.NodeDetails(r.Context(), nodeID, f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleMapEdge(w http.ResponseWriter, r *http.Request) {
	edgeID := chi.URLParam(r, "id")
	edgeID = strings.ReplaceAll(edgeID, "%3E", ">")
	f := parseFilter(r)
	data, err := s.store.EdgeDetails(r.Context(), edgeID, f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	if strings.TrimSpace(f.Search) == "" {
		writeJSON(w, http.StatusOK, map[string]any{"data": []model.TopItem{}})
		return
	}
	data, err := s.store.Search(r.Context(), f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": data})
}

func (s *Server) handleInventoryImport(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Items []model.InventoryAsset `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	for i := range body.Items {
		if strings.TrimSpace(body.Items[i].AssetID) == "" {
			body.Items[i].AssetID = strconv.FormatInt(time.Now().UTC().UnixNano()+int64(i), 36)
		}
	}
	if err := s.store.ImportInventory(r.Context(), body.Items); err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"imported": len(body.Items)})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{"error": err.Error()})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
