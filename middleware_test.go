package pine

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddingMiddleware(t *testing.T) {
	p := New()
	if len(p.Middlewares) != 0 {
		t.Errorf("middlewares was not expected length %d, instead was %d", 0, len(p.Middlewares))
	}

	p.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})

	if len(p.Middlewares) != 1 {
		t.Errorf("middlewares was not expected length %d, instead was %d", 1, len(p.Middlewares))
	}
}

func TestMiddlewareEarlyReturn(t *testing.T) {
	p := New()
	p.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				w.Write([]byte("hello from middleware"))
			} else {
				next.ServeHTTP(w, r)
			}
		})
	})
	p.Handle("/testing", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from handler"))
	})

	req := httptest.NewRequest("GET", "https://lukaswerner.com/", nil)
	resp := httptest.NewRecorder()
	p.ServeHTTP(resp, req)

	result := resp.Result()
	if b, _ := io.ReadAll(result.Body); bytes.Compare(b, []byte("hello from middleware")) != 0 {
		t.Errorf("the body response did not match the expected body: %s != %s", string(b), "hello from middleware")
	}

	req = httptest.NewRequest("GET", "https://lukaswerner.com/testing", nil)
	resp = httptest.NewRecorder()
	p.ServeHTTP(resp, req)

	result = resp.Result()
	if b, _ := io.ReadAll(result.Body); bytes.Compare(b, []byte("hello from handler")) != 0 {
		t.Errorf("the body response did not match the expected body: %s != %s", string(b), "hello from handler")
	}

}
