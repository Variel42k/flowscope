package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/flowscope/flowscope/internal/auth"
)

func (s *Server) handleAuthMe(w http.ResponseWriter, r *http.Request) {
	user := auth.AuthUserFromContext(r.Context())
	if strings.TrimSpace(user.Username) == "" {
		writeErr(w, http.StatusUnauthorized, errUnauthorized())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user": user.Username,
		"role": user.Role,
	})
}

func (s *Server) handleOIDCStart(w http.ResponseWriter, r *http.Request) {
	if s.oidc == nil || !s.oidc.Enabled() {
		writeErr(w, http.StatusBadRequest, errOIDCDisabled())
		return
	}
	authURL, state, err := s.oidc.StartAuthURL(r.Context())
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"authorize_url": authURL,
		"state":         state,
	})
}

func (s *Server) handleOIDCCallback(w http.ResponseWriter, r *http.Request) {
	if s.oidc == nil || !s.oidc.Enabled() {
		writeErr(w, http.StatusBadRequest, errOIDCDisabled())
		return
	}
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	username, role, err := s.oidc.ExchangeCode(r.Context(), state, code)
	if err != nil {
		s.redirectOIDCResult(w, r, "", "", err)
		return
	}
	token, err := s.auth.Sign(username, role)
	if err != nil {
		s.redirectOIDCResult(w, r, "", "", err)
		return
	}
	s.redirectOIDCResult(w, r, token, username, nil)
}

func (s *Server) redirectOIDCResult(w http.ResponseWriter, r *http.Request, token string, username string, resultErr error) {
	target := strings.TrimSpace(s.cfg.OIDCSuccessRedirect)
	if target == "" {
		target = "http://localhost:5173/oidc/callback"
	}
	u, err := url.Parse(target)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	q := u.Query()
	if resultErr != nil {
		q.Set("error", resultErr.Error())
	} else {
		q.Set("token", token)
		q.Set("user", username)
		q.Set("role", "")
		claims, parseErr := s.auth.Parse(token)
		if parseErr == nil {
			q.Set("role", claims.Role)
		}
		q.Set("auth_mode", "oidc")
	}
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func errOIDCDisabled() error {
	return errBadRequest("oidc is not enabled")
}

func errUnauthorized() error {
	return errBadRequest("unauthorized")
}

func errBadRequest(msg string) error {
	return &httpError{Message: msg}
}

type httpError struct {
	Message string
}

func (e *httpError) Error() string {
	return e.Message
}
