package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/lukasmwerner/pine"
)

type responseWriterBuffer struct {
	bodyBuffer *bytes.Buffer
	headers    http.Header
	statusCode int
}

func (w *responseWriterBuffer) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = 200
	}
	return w.bodyBuffer.Write(b)
}
func (w *responseWriterBuffer) Header() http.Header {
	return w.headers
}
func (w *responseWriterBuffer) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (wb *responseWriterBuffer) Copy(w http.ResponseWriter) {
	// Copy all the headers
	for k, vs := range wb.Header() {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(wb.statusCode)
	io.Copy(w, wb.bodyBuffer)
}

func STolinskiTiming(slow time.Duration, middle time.Duration) pine.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			rspWB := &responseWriterBuffer{
				bodyBuffer: bytes.NewBuffer(nil),
			}

			start := time.Now()
			next.ServeHTTP(rspWB, r)
			finish := time.Now()
			duration := finish.Sub(start)
			w.Header().Set("X-Duration", "rocket!")
			if duration.Microseconds() >= slow.Microseconds() {
				w.Header().Set("X-Duration", "turtle...")
			}
			if duration.Microseconds() < slow.Microseconds() && duration.Microseconds() >= middle.Microseconds() {
				w.Header().Set("X-Duration", "rabbit..")
			}
			rspWB.Copy(w)
		})
	}
}
