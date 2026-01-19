package api

import (
	"net/http"

	"control-plane/internal/api/middleware"
)

func Router() http.Handler {
	mux := http.NewServeMux()
	return middleware.Auth(mux)
}
