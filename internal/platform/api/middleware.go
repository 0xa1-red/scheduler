package api

import (
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *statusWriter) WriteHeader(s int) {
	w.status = s
	w.ResponseWriter.WriteHeader(s)
}

func (w *statusWriter) Write(d []byte) (int, error) {
	n, err := w.ResponseWriter.Write(d)
	w.size += n
	return n, err
}

type LoggingMiddleware struct{}

func (lmw *LoggingMiddleware) MiddlewareFn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := statusWriter{
			ResponseWriter: w,
			status:         200,
			size:           0,
		}
		start := time.Now()
		next.ServeHTTP(&sw, r)
		t := time.Since(start)
		logger.Infow("request", "method", r.Method, "url", r.URL.String(), "status", sw.status, "response_size", sw.size, "response_time", t.String())
	})
}
