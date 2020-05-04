package echo

import "net/http"

type (
	Group struct {
		common
		host       string
		prefix     string
		middleware []MiddlewareFunc
		echo       *Echo
	}
)

func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
	if len(g.middleware) == 0 {
		return
	}
	g.Any("", NotFoundHandler)
	g.Any("/*", NotFoundHandler)
}

func (g *Group) CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodConnect, path, h, m...)
}

func (g *Group) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodDelete, path, h, m...)
}

func (g *Group) GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodGet, path, h, m...)
}

func (g *Group) HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodHead, path, h, m...)
}

func (g *Group) OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodOptions, path, h, m...)
}

func (g *Group) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPatch, path, h, m...)
}

func (g *Group) POST(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPost, path, h, m...)
}

func (g *Group) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPut, path, h, m...)
}

func (g *Group) TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodTrace, path, h, m...)
}

func (g *Group) Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route {
	routes := make([]*Route, len(methods))
	for i, m := range methods {
		routes[i] = g.Add(m, path, handler, middleware...)
	}
	return routes
}

func (g *Group) Match(methods []string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route {
	routes := make([]*Route, len(methods))
	for i, m := range methods {
		routes[i] = g.Add(m, path, handler, middleware...)
	}
	return routes
}

func (g *Group) Group(prefix string, middlewareFunc ...MiddlewareFunc) {
	m := make([]MiddlewareFunc, 0, len(g.middleware)+len(middlewareFunc))
	m = append(m, g.middleware)
	m = append(m, middlewareFunc)
	sg = g.echo.Group(g.prefix+prefix, m...)
	sg.host = g.host
	return
}

func (g *Group) Static(prefix, root string) {
	g.static(prefix, root, g.GET)
}

func (g *Group) File(path, file string) {
	g.file(path, file, g.GET)
}

func (g *Group) Add(method, path string, handler HandlerFunc, middlewareFunc ...MiddlewareFunc) *Route {
	m := make([]MiddlewareFunc, 0, len(g.middleware)+len(middlewareFunc))
	m = append(m, g.middleware...)
	m = append(m, middlewareFunc...)
	return g.echo.add(g.host, method, g.prefix+path, handler, m...)
}
