package runtime

import "net/http"

func Router() http.Handler {
	mux := http.NewServeMux()
	return Auth(mux)
}
