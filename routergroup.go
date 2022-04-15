package gem

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	engie    *Engine
	basePath string
	root     bool
	handlers []HandlerFunc
}

func (group *RouterGroup) Use(handlers ...HandlerFunc) {
	group.handlers = append(group.handlers, handlers...)
}

func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		engie:    group.engie,
		basePath: group.combineAbsoltePath(relativePath),
		handlers: group.combineHandlers(handlers),
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	pattern := path.Join(relativePath, "/*filepath")
	group.GET(pattern, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.basePath, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.SetStatusCode(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (group *RouterGroup) addRoute(method string, relativePath string, handlers ...HandlerFunc) {
	pattern := group.combineAbsoltePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engie.addRoute(method, pattern, handlers)
}

func (group *RouterGroup) GET(relativePath string, handler HandlerFunc) {
	group.addRoute(http.MethodGet, relativePath, handler)
}

func (group *RouterGroup) POST(relativePath string, handler HandlerFunc) {
	group.addRoute(http.MethodPost, relativePath, handler)
}

func (group *RouterGroup) PUT(relativePath string, handler HandlerFunc) {
	group.addRoute(http.MethodPut, relativePath, handler)
}

func (group *RouterGroup) DELETE(relativePath string, handler HandlerFunc) {
	group.addRoute(http.MethodDelete, relativePath, handler)
}

func (group *RouterGroup) combineAbsoltePath(relativePath string) string {
	if relativePath == "" {
		return group.basePath
	}
	return path.Join(group.basePath, relativePath)
}

func (group *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	size := len(group.handlers) + len(handlers)
	if size >= 9 {
		panic("")
	}
	mergeHandlers := make([]HandlerFunc, size)
	copy(mergeHandlers, group.handlers)
	copy(mergeHandlers[len(group.handlers):], handlers)
	return mergeHandlers
}
