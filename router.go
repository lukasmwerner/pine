package pine

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

var allMethods = []string{"GET", "POST", "DELETE", "PUT", "HEAD", "PATCH", "CONNECT", "OPTIONS", "TRACE"}
var variableRegex = regexp.MustCompile(`(?m)\{([a-zA-Z]*)\}`)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("not found."))
}

func New() *router {
	return &router{
		RootNode: &node{
			Path:     "", // this should match anything
			children: []*node{},
			Name:     "root",
		},
		NotFoundHandler: http.HandlerFunc(notFoundHandler),
		Middlewares:     []MiddlewareFunc{},
	}
}

func Var(r *http.Request, k string) string {
	return r.Context().Value("pine:" + k).(string)
}

type MiddlewareFunc func(http.Handler) http.Handler

type router struct {
	RootNode        *node
	NotFoundHandler http.Handler
	Middlewares     []MiddlewareFunc
}

func (r *router) Handle(tpl string, handler http.HandlerFunc) {
	path_parts := strings.Split(tpl, "/")[1:]

	// 1. find the farthest node
	var farthest *node = r.RootNode
	farthest, count := findMatchingPathNode(r.RootNode, path_parts, 0)

	// 2. insert as much as needed
	path_parts_to_insert := path_parts[count:]
	for len(path_parts_to_insert) > 0 {
		path_part := path_parts_to_insert[0]
		new_node := &node{
			children: []*node{},
			Path:     path_part,
		}
		if variableRegex.MatchString(path_part) {
			new_node.Path = ""
			new_node.variable = variableRegex.FindStringSubmatch(path_part)[1]
		}

		farthest.children = append(farthest.children, new_node)
		path_parts_to_insert = path_parts_to_insert[1:] // pop off the left
		farthest = new_node
	}

	// 3. insert all methods for leaves
	for _, method := range allMethods {
		new_node := &node{
			Method:   method,
			children: []*node{},
			handler:  http.HandlerFunc(handler),
		}
		farthest.children = append(farthest.children, new_node)
	}
}

func (r *router) Use(middleware MiddlewareFunc) {
	r.Middlewares = append(r.Middlewares, middleware)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := findMatchingHandler(r.RootNode, req)
	if handler == nil {
		handler = r.NotFoundHandler
	}
	varsMap := makeVarsFromRequest(r.RootNode, req)
	ctx := req.Context()
	for k, v := range varsMap {
		ctx = context.WithValue(ctx, "pine:"+k, v)
	}

	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		handler = r.Middlewares[i](handler)
	}
	handler.ServeHTTP(w, req.WithContext(ctx))
}

type node struct {
	children []*node
	handler  http.Handler
	Path     string
	Method   string
	Name     string
	variable string
}

func findMatchingHandler(root *node, r *http.Request) http.Handler {
	request_path := r.URL.Path
	// get the folders and skip the first element as that is an empty string
	request_path_parts := strings.Split(request_path, "/")[1:]
	branch, _ := findMatchingPathNode(root, request_path_parts, 0)
	leaf := findMatchingMethodNode(branch, r.Method)
	return leaf.handler
}

func findMatchingPathNode(root *node, request_path []string, count int) (*node, int) {
	if len(request_path) == 0 {
		return root, count + 1
	}
	for _, child := range root.children {
		if child.Path == request_path[0] {
			if child.Method != "" {
				return root, count + 1
			}
			return findMatchingPathNode(child, request_path[1:], count+1)
		}
		if child.Path == "" {
			if child.Method != "" {
				return root, count + 1
			}
			return findMatchingPathNode(child, request_path[1:], count+1)
		}
	}
	// not sure about the root case maybe just return the root node?
	return root, count
}
func findMatchingMethodNode(root *node, method string) *node {
	for _, child := range root.children {
		if child.Method == method {
			return child
		}
	}
	return root
}
func makeVarsFromRequest(start *node, r *http.Request) map[string]string {
	request_path := r.URL.Path
	vars := make(map[string]string)
	request_path_parts := strings.Split(request_path, "/")[1:]

	// loop over children and fill up the vars as we encounter them
	var root *node = start
	for {
		if len(request_path_parts) == 0 {
			break
		}
		for _, child := range root.children {
			if child.Path == request_path_parts[0] {
				root = child
				break
			}
			if child.Path == "" && child.variable != "" {
				vars[child.variable] = request_path_parts[0]
			}
		}
		request_path_parts = request_path_parts[1:]
	}

	return vars
}
