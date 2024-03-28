//go:build !solution

package requestlog

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"runtime"
	"time"

	"go.uber.org/zap"
)

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func GenerateRequestID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)

	if err != nil {
		return time.Now().Format(time.RFC3339Nano)
	}

	return hex.EncodeToString(bytes)
}

func Log(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := GenerateRequestID()
			LogInit(l, requestID, r)
			start := time.Now()
			wrapper := &ResponseWriter{ResponseWriter: w, status: http.StatusOK}

			defer func() {
				duration := time.Since(start)

				if err := recover(); err != nil {
					bytes := make([]byte, 4096)
					bytes = bytes[:runtime.Stack(bytes, false)]
					LogError(l, requestID, duration, wrapper, r, bytes)
					panic(err)
				}

				LogInfo(l, requestID, duration, wrapper, r)
			}()

			next.ServeHTTP(wrapper, r)
		})
	}
}

func LogInfo(l *zap.Logger, requestID string, duration time.Duration, wrapper *ResponseWriter, r *http.Request) {
	l.Info("request finished",
		zap.String("request_id", requestID),
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
		zap.Duration("duration", duration),
		zap.Int("status_code", wrapper.status),
	)
}

func LogError(l *zap.Logger, requestID string, duration time.Duration, wrapper *ResponseWriter, r *http.Request, bytes []byte) {
	l.Error("request panicked",
		zap.String("request_id", requestID),
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
		zap.Duration("duration", duration),
		zap.Int("status_code", wrapper.status),
		zap.String("stack_trace", string(bytes)),
	)
}

func LogInit(l *zap.Logger, requestID string, r *http.Request) {
	l.Info("request started",
		zap.String("request_id", requestID),
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)
}
