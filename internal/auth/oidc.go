package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCConfig struct {
	Enabled            bool
	IssuerURL          string
	ClientID           string
	ClientSecret       string
	RedirectURL        string
	SuccessRedirectURL string
	Scopes             []string
	RoleClaim          string
	AdminUsers         []string
	AllowedDomains     []string
	AdminRole          string
	ViewerRole         string
}

type OIDCProvider struct {
	cfg      OIDCConfig
	initOnce sync.Once
	initErr  error

	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2   oauth2.Config

	stateMu sync.Mutex
	states  map[string]oidcState
}

type oidcState struct {
	Nonce     string
	ExpiresAt time.Time
}

func NewOIDCProvider(cfg OIDCConfig) *OIDCProvider {
	if cfg.AdminRole == "" {
		cfg.AdminRole = "admin"
	}
	if cfg.ViewerRole == "" {
		cfg.ViewerRole = "viewer"
	}
	if len(cfg.Scopes) == 0 {
		cfg.Scopes = []string{"openid", "profile", "email"}
	}
	if cfg.RoleClaim == "" {
		cfg.RoleClaim = "role"
	}
	return &OIDCProvider{
		cfg:    cfg,
		states: make(map[string]oidcState),
	}
}

func (p *OIDCProvider) Enabled() bool {
	return p != nil && p.cfg.Enabled
}

func (p *OIDCProvider) ensureInit(ctx context.Context) error {
	if !p.Enabled() {
		return errors.New("oidc is disabled")
	}
	p.initOnce.Do(func() {
		if strings.TrimSpace(p.cfg.IssuerURL) == "" || strings.TrimSpace(p.cfg.ClientID) == "" || strings.TrimSpace(p.cfg.RedirectURL) == "" {
			p.initErr = errors.New("oidc config is incomplete")
			return
		}
		provider, err := oidc.NewProvider(ctx, p.cfg.IssuerURL)
		if err != nil {
			p.initErr = err
			return
		}
		p.provider = provider
		p.verifier = provider.Verifier(&oidc.Config{ClientID: p.cfg.ClientID})
		p.oauth2 = oauth2.Config{
			ClientID:     p.cfg.ClientID,
			ClientSecret: p.cfg.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  p.cfg.RedirectURL,
			Scopes:       p.cfg.Scopes,
		}
	})
	return p.initErr
}

func (p *OIDCProvider) StartAuthURL(ctx context.Context) (string, string, error) {
	if err := p.ensureInit(ctx); err != nil {
		return "", "", err
	}
	state, err := randomToken(24)
	if err != nil {
		return "", "", err
	}
	nonce, err := randomToken(24)
	if err != nil {
		return "", "", err
	}
	p.stateMu.Lock()
	p.states[state] = oidcState{
		Nonce:     nonce,
		ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	}
	p.cleanupStatesLocked()
	p.stateMu.Unlock()
	url := p.oauth2.AuthCodeURL(state, oauth2.SetAuthURLParam("nonce", nonce))
	return url, state, nil
}

func (p *OIDCProvider) ExchangeCode(ctx context.Context, state, code string) (string, string, error) {
	if err := p.ensureInit(ctx); err != nil {
		return "", "", err
	}
	state = strings.TrimSpace(state)
	code = strings.TrimSpace(code)
	if state == "" || code == "" {
		return "", "", errors.New("missing state or code")
	}
	p.stateMu.Lock()
	st, ok := p.states[state]
	delete(p.states, state)
	p.stateMu.Unlock()
	if !ok || time.Now().UTC().After(st.ExpiresAt) {
		return "", "", errors.New("invalid or expired oidc state")
	}
	token, err := p.oauth2.Exchange(ctx, code)
	if err != nil {
		return "", "", err
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || strings.TrimSpace(rawIDToken) == "" {
		return "", "", errors.New("id_token is missing")
	}
	idToken, err := p.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", "", err
	}
	var claims map[string]any
	if err := idToken.Claims(&claims); err != nil {
		return "", "", err
	}
	if n, ok := claims["nonce"].(string); ok && n != "" && n != st.Nonce {
		return "", "", errors.New("oidc nonce mismatch")
	}
	username := pickUsername(claims)
	if username == "" {
		return "", "", errors.New("unable to determine user identity from oidc claims")
	}
	if !p.allowedByDomain(username) {
		return "", "", fmt.Errorf("user domain is not allowed: %s", username)
	}
	role := p.resolveRole(username, claims)
	return username, role, nil
}

func (p *OIDCProvider) resolveRole(username string, claims map[string]any) string {
	adminRole := p.cfg.AdminRole
	viewerRole := p.cfg.ViewerRole
	if roleClaimVal, ok := claims[p.cfg.RoleClaim]; ok {
		if claimRole := parseRoleValue(roleClaimVal); claimRole != "" {
			claimRole = strings.ToLower(strings.TrimSpace(claimRole))
			if claimRole == strings.ToLower(adminRole) {
				return adminRole
			}
			if claimRole == strings.ToLower(viewerRole) {
				return viewerRole
			}
		}
	}
	lowerUser := strings.ToLower(strings.TrimSpace(username))
	for _, admin := range p.cfg.AdminUsers {
		if lowerUser == strings.ToLower(strings.TrimSpace(admin)) {
			return adminRole
		}
	}
	return viewerRole
}

func (p *OIDCProvider) allowedByDomain(username string) bool {
	if len(p.cfg.AllowedDomains) == 0 {
		return true
	}
	parts := strings.Split(username, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.ToLower(strings.TrimSpace(parts[1]))
	for _, allowed := range p.cfg.AllowedDomains {
		if domain == strings.ToLower(strings.TrimSpace(allowed)) {
			return true
		}
	}
	return false
}

func (p *OIDCProvider) cleanupStatesLocked() {
	now := time.Now().UTC()
	for state, st := range p.states {
		if now.After(st.ExpiresAt) {
			delete(p.states, state)
		}
	}
}

func pickUsername(claims map[string]any) string {
	for _, key := range []string{"preferred_username", "email", "upn", "sub"} {
		if v, ok := claims[key].(string); ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func parseRoleValue(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []any:
		for _, item := range x {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				return s
			}
		}
	}
	return ""
}

func randomToken(n int) (string, error) {
	if n < 16 {
		n = 16
	}
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
