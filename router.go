package gem

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*node
}

func newRoter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	trie, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := trie.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *Router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if params != nil {
		c.Params = params
	}
	if node != nil {
		key := c.Method + "-" + node.pattern
		if handler, ok := r.handlers[key]; ok {
			c.handlers = append(c.handlers, handler)
		}
	} else {
		http.Error(c.Writer, "404 NOT FOUND: %s\n", http.StatusInternalServerError)
	}
	c.Next()
}

func parsePattern(pattern string) []string {
	values := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range values {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
