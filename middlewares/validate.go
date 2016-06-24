package middlewares

import "net/http"

func ValidateJsonHeader(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			http.Error(w, "header Accept is not set to application/json", http.StatusBadRequest)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "header Content-Type is not set to application/json", http.StatusBadRequest)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
