package middlewares

import (
	"log"
	"net/http"

	"github.com/lukasmwerner/pine"
)

type LoggerConfig struct {
	l *log.Logger
}

// HTTPLogger Logs the HTTP Request/Responses in the following
//
// format: 2023/04/02 12:48:37 host: localhost:8080 method: GET uri: / status: 200 Ok
func HTTPLogger() pine.MiddlewareFunc {
	c := LoggerConfig{}
	c.l = log.Default()
	c.l.SetFlags(log.Ldate | log.Ltime)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := &loggingWriter{w, 200}
			next.ServeHTTP(lw, r)
			c.l.Printf(
				"host: %s method: %s uri: %s status: %d %s\n",
				r.Host,
				r.Method,
				r.URL.RequestURI(),
				lw.statusCode,
				http.StatusText(lw.statusCode),
			)
		})
	}
}

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

func (l *loggingWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}
