package gem

import (
	"net/http"
	"sync"
	"text/template"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	pool       sync.Pool
	trees      methodTrees
	htmlRender *template.Template
	funcMap    template.FuncMap
}

func Default() *Engine {
	r := New()
	r.Use(Logger(), Recovery())
	return r
}

func New() *Engine {
	engie := &Engine{
		RouterGroup: &RouterGroup{
			handlers: nil,
			root:     true,
			basePath: "/",
		},
		trees: make(methodTrees, 0, 9),
	}
	engie.RouterGroup.engie = engie
	engie.pool.New = func() interface{} {
		return engie.allocateContext()
	}
	return engie
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlRender = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) allocateContext() *Context {
	return &Context{engine: e}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := e.pool.Get().(*Context)
	c.Request = r
	c.Writer = w
	c.reset()

	e.handleHTTPRequest(c)

	e.pool.Put(c)
}

func (e *Engine) addRoute(method string, pattern string, handlers []HandlerFunc) {
	parts := parsePattern(pattern)
	root := e.trees.get(method)
	if root == nil {
		root = new(node)
		root.pattern = "/"
		e.trees = append(e.trees, methodTree{method: method, root: root})
	}
	root.insert(pattern, parts, 0, handlers)
}

func (e *Engine) handleHTTPRequest(c *Context) {
	method := c.Request.Method
	path := c.Request.URL.Path
	root := e.trees.get(method)
	node, params := root.getValue(path)
	if params != nil {
		c.Params = params
	}
	if node != nil {
		c.Method = method
		c.Path = node.pattern
		c.handlers = node.handlers
		c.Next()
	} else {
		serveError(c, http.StatusNotFound, []byte("404 NOT FOUND\n"))
	}
}

func serveError(c *Context, code int, defaultMessage []byte) {
	c.SetStatusCode(code)
	c.Next()
	c.Writer.Write(defaultMessage)
}
