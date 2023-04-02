package pine

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNodeTraversal(t *testing.T) {
	root := &node{
		Path:     "", // this should match anything
		children: []*node{},
		Name:     "root",
	}
	root.children = append(root.children, &node{
		Path: "hello",
		Name: "hello folder",
		children: []*node{{
			Path: "", // this should match anything
			Name: "hello name handler",
			children: []*node{{
				Method: "GET",
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("welcome to the forest"))
				}),
				Name: "get handler",
			}, {
				Method: "POST",
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("thank you for the mail"))
				}),
				Name: "post handler",
			}},
		}, {
			Method: "GET",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("hello there!"))
			}),
		}},
	})

	req := httptest.NewRequest("GET", "https://lukaswerner.com/hello/world", nil)

	handler := findMatchingHandler(root, req)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if b, _ := io.ReadAll(resp.Body); bytes.Compare(b, []byte("welcome to the forest")) != 0 {
		t.Errorf("the body response did not match the expected body: %s != %s", string(b), "welcome to the forest")
	}

	req = httptest.NewRequest("POST", "https://lukaswerner.com/hello/world", nil)

	handler = findMatchingHandler(root, req)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp = w.Result()
	if b, _ := io.ReadAll(resp.Body); bytes.Compare(b, []byte("thank you for the mail")) != 0 {
		t.Errorf("the body response did not match the expected body: %s != %s", string(b), "thank you for the mail")
	}

	req = httptest.NewRequest("GET", "https://lukaswerner.com/hello", nil)

	handler = findMatchingHandler(root, req)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp = w.Result()
	if b, _ := io.ReadAll(resp.Body); bytes.Compare(b, []byte("hello there!")) != 0 {
		t.Errorf("the body response did not match the expected body: %s != %s", string(b), "hello there!")
	}
}

func TestInsertingNodes(t *testing.T) {
	p := New()
	p.Handle("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})

	req := httptest.NewRequest("GET", "https://lukaswerner.com/hello/world", nil)
	handler := findMatchingHandler(p.RootNode, req)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if b, _ := io.ReadAll(resp.Body); bytes.Compare(b, []byte("hello world!")) != 0 {
		t.Errorf("the body response did not match the expected body: %q != %q", string(b), "hello world!")
	}
}

func TestPathVarsInsertion(t *testing.T) {
	p := New()
	p.Handle("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})

	req := httptest.NewRequest("GET", "https://lukaswerner.com/hello/world", nil)
	vars := makeVarsFromRequest(p.RootNode, req)
	if len(vars) > 1 || len(vars) <= 0 {
		t.Errorf("vars was %d not the correct size as expected %d", len(vars), 1)
	}
	if vars["name"] != "world" {
		t.Errorf("value for \"name\" was not what was expected: %q got: %q", "world", vars["name"])
	}
}

func TestPathVarsUsage(t *testing.T) {
	value := ""
	p := New()
	p.Handle("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		value = Var(r, "name")
		w.Write([]byte("hello " + value + "!"))
	})

	req := httptest.NewRequest("GET", "https://lukaswerner.com/hello/world", nil)
	w := httptest.NewRecorder()
	p.ServeHTTP(w, req)

	resp := w.Result()
	if b, _ := io.ReadAll(resp.Body); bytes.Compare(b, []byte("hello world!")) != 0 {
		t.Errorf("the body response did not match the expected body: %s != %s", string(b), "hello world!")
	}

	if value != "world" {
		t.Errorf("answer for path variable loading was not what was expected: %s but got %s instead", "world", value)
	}
}
