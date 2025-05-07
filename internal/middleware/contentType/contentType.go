package contentType

import (
	"net/http"
	"strings"
)

var jsonPrefixes = []string{
	"/auth",
	"/user",
	"/server",
	"/forecast",
	"/metrics",
}

func isJsonRoute(path string) bool {
	for _, prefix := range jsonPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func JsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isJsonRoute(r.URL.Path) && r.Header.Get("Upgrade") != "websocket" {
			w.Header().Set("Content-Type", "application/json")
		}

		next.ServeHTTP(w, r)
	})
}
