package gem

import (
	"net/http"
	"strings"
	"text/template"
)

type Engine struct {
	*RouterGroup
	router     *Router
	groups     []*RouterGroup
	htmlRender *template.Template
	funcMap    template.FuncMap
}

func New() *Engine {
	r := &Engine{router: newRoter()}
	r.RouterGroup = &RouterGroup{engie: r}
	r.groups = []*RouterGroup{r.RouterGroup}
	r.Use(Logger(), Recovery())
	return r
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

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRoute("POST", pattern, handler)
}

func (e *Engine) PUT(pattern string, handler HandlerFunc) {
	e.router.addRoute("PUT", pattern, handler)
}

func (e *Engine) DELETE(pattern string, handler HandlerFunc) {
	e.router.addRoute("DELETE", pattern, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middlewares := []HandlerFunc{}
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.basePath) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	ctx := newContext(w, r)
	ctx.handlers = middlewares
	ctx.engine = e
	e.router.handle(ctx)
}
