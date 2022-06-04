package routes

import (
	"context"
	"net/http"
	"otus_dialog_go/internal/constants"
	"otus_dialog_go/internal/otusprometheus"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		otusprometheus.IncRequestCount()
		startTime := time.Now().UnixNano() / 1000000
		requestID := r.Header.Get("X-REQUEST-ID")
		ctx := context.WithValue(r.Context(), constants.RequestIDKey, requestID)
		lw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lw, r.WithContext(ctx))
		endTime := time.Now().UnixNano() / 1000000
		requestTime := endTime - startTime
		otusprometheus.AddRequestTimeData(float64(requestTime))
		if lw.statusCode >= 500 {
			otusprometheus.IncErrorsCount()
		}
	})
}
