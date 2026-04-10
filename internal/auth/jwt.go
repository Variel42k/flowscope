package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret []byte
}

func NewManager(secret string) *Manager {
	return &Manager{secret: []byte(secret)}
}

func (m *Manager) Sign(username string) (string, error) {
	claims := Claims{
		Username: username,
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
	if v := ctx.Value(userContextKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func Middleware(manager *Manager, authDisabled bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authDisabled {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, "dev-admin")))
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
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, claims.Username)))
		})
	}
}
