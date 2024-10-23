package mid

import (
	"fmt"
	"net/http"
	"time"

	"github.com/3bd-dev/wallet-service/pkg/logger"
)

// Logger writes information about the request to the logs.
func Logger(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info(r.Context(), "request started", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr)

			// Create a response writer that captures the status code
			ww := &rWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(ww, r)

			log.Info(r.Context(), "request completed", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr,
				"statuscode", ww.statusCode, "since", time.Since(now).String())
		})
	}
}

type rWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *rWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
