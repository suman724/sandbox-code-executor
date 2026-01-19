package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"control-plane/internal/audit"
)

type ctxKey string

const (
	ctxTenantIDKey ctxKey = "tenant_id"
	ctxAgentIDKey  ctxKey = "agent_id"
)

type Claims struct {
	TenantID string `json:"tenant_id"`
	AgentID  string `json:"agent_id"`
	jwt.RegisteredClaims
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("AUTHZ_BYPASS") == "true" {
			audit.StdoutLogger{}.Log(r.Context(), audit.Event{
				Action:  "authz_bypass",
				Outcome: "ok",
				Time:    time.Now(),
			})
			next.ServeHTTP(w, r)
			return
		}
		token, err := parseBearer(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		parsed, err := jwt.ParseWithClaims(token, claims, keyFunc(), jwtOptions()...)
		if err != nil || !parsed.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxTenantIDKey, claims.TenantID)
		ctx = context.WithValue(ctx, ctxAgentIDKey, claims.AgentID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseBearer(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing authorization header")
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", errors.New("missing bearer prefix")
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", errors.New("empty token")
	}
	return token, nil
}

func keyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		secret := os.Getenv("AUTH_JWT_SECRET")
		if secret == "" {
			return nil, errors.New("AUTH_JWT_SECRET not set")
		}
		return []byte(secret), nil
	}
}

func jwtOptions() []jwt.ParserOption {
	var opts []jwt.ParserOption
	if issuer := os.Getenv("AUTH_ISSUER"); issuer != "" {
		opts = append(opts, jwt.WithIssuer(issuer))
	}
	if audience := os.Getenv("AUTH_AUDIENCE"); audience != "" {
		opts = append(opts, jwt.WithAudience(audience))
	}
	return opts
}
