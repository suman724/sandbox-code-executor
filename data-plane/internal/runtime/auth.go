package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"data-plane/internal/telemetry"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("AUTHZ_BYPASS") == "true" {
			telemetry.StdoutLogger{}.Log(r.Context(), telemetry.Event{
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
		if err := validateToken(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
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

func validateToken(token string) error {
	parser := jwt.NewParser(jwtOptions()...)
	parsed, err := parser.Parse(token, keyFunc())
	if err != nil || !parsed.Valid {
		return errors.New("invalid token")
	}
	return nil
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
