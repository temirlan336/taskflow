package middleware

import (
	"context"
	"net/http"
	"os"
)

const apiKeyContextKey string = "apiKey"

func AuthMiddleware(next http.Handler) http.Handler {
	expectedKey := os.Getenv("API_KEY")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" || apiKey != expectedKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), apiKeyContextKey, apiKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAPIKeyFromContext(ctx context.Context) (string, bool) {
	apiKey, ok := ctx.Value(apiKeyContextKey).(string)
	return apiKey, ok
}
