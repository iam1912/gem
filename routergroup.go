package gem

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	engie       *Engine
	middlewares []HandlerFunc
	basePath    string
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) Group(basePath string) *RouterGroup {
	newgroup := &RouterGroup{
		engie:    group.engie,
		basePath: group.basePath + basePath,
	}
	group.engie.groups = append(group.engie.groups, newgroup)
	return newgroup
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
		c.StatusCode = 200
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) addRoute(method string, relativePath string, handler HandlerFunc) {
	pattern := group.basePath + relativePath
	group.engie.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(relativePath string, handler HandlerFunc) {
	group.addRoute("GET", relativePath, handler)
}

func (group *RouterGroup) POST(relativePath string, handler HandlerFunc) {
	group.addRoute("POST", relativePath, handler)
}

func (group *RouterGroup) PUT(relativePath string, handler HandlerFunc) {
	group.addRoute("PUT", relativePath, handler)
}

func (group *RouterGroup) DELETE(relativePath string, handler HandlerFunc) {
	group.addRoute("DELETE", relativePath, handler)
}
