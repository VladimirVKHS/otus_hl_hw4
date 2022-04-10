package routes

import (
	"context"
	"net/http"
	"otus_dialog_go/internal/constants"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-REQUEST-ID")
		ctx := context.WithValue(r.Context(), constants.RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
