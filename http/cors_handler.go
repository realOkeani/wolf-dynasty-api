package http

import (
"net/http"
)

func CorsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("access-control-allow-origin", "*")
			if r.Method == http.MethodOptions {
				w.Header().Set("access-control-allow-headers", "authorization")
				w.Header().Set("access-control-allow-methods", "PATCH,PUT,POST,OPTIONS,GET,DELETE")
				return
			}
			h.ServeHTTP(w, r)
		})
}