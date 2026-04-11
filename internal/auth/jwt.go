package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/flowscope/flowscope/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret []byte
}

func NewManager(secret string) *Manager {
	return &Manager{secret: []byte(secret)}
}

func (m *Manager) Sign(username string, role string) (string, error) {
	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "flowscope",
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(_ *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

type contextKey string

const userContextKey contextKey = "user"

func UserFromContext(ctx context.Context) string {
	return AuthUserFromContext(ctx).Username
}

func UserRoleFromContext(ctx context.Context) string {
	return AuthUserFromContext(ctx).Role
}

func AuthUserFromContext(ctx context.Context) model.AuthUser {
	if v := ctx.Value(userContextKey); v != nil {
		if user, ok := v.(model.AuthUser); ok {
			return user
		}
		if s, ok := v.(string); ok {
			return model.AuthUser{Username: s, Role: "admin"}
		}
	}
	return model.AuthUser{}
}

func Middleware(manager *Manager, authDisabled bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authDisabled {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, model.AuthUser{
					Username: "dev-admin",
					Role:     "admin",
				})))
				return
			}
			hdr := r.Header.Get("Authorization")
			if !strings.HasPrefix(strings.ToLower(hdr), "bearer ") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			token := strings.TrimSpace(hdr[7:])
			claims, err := manager.Parse(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, model.AuthUser{
				Username: claims.Username,
				Role:     claims.Role,
			})))
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	required := strings.ToLower(strings.TrimSpace(role))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := AuthUserFromContext(r.Context())
			if required == "" || strings.EqualFold(user.Role, required) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		role = strings.ToLower(strings.TrimSpace(role))
		if role != "" {
			allowed[role] = struct{}{}
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(allowed) == 0 {
				next.ServeHTTP(w, r)
				return
			}
			user := AuthUserFromContext(r.Context())
			if _, ok := allowed[strings.ToLower(strings.TrimSpace(user.Role))]; ok {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}
