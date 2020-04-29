package echo

import "net/http"

type (
	Router struct {
		tree   *node
		routes map[string]*Route
		echo   *Echo
	}
	node struct {
		kind          kind
		label         byte
		prefix        string
		parent        *node
		children      children
		ppath         string
		pnames        []string
		methodHandler *methodHandler
	}
	kind          uint8
	children      []*node
	methodHandler struct {
		connect  HandlerFunc
		delete   HandlerFunc
		get      HandlerFunc
		head     HandlerFunc
		options  HandlerFunc
		patch    HandlerFunc
		post     HandlerFunc
		propfind HandlerFunc
		put      HandlerFunc
		trace    HandlerFunc
		report   HandlerFunc
	}
)

const (
	skind kind = iota
	pkind
	akind
)

func NewRouter(e *Echo) *Router {
	return &Router{
		tree: &node{
			methodHandler: new(methodHandler),
		},
		routes: map[string]*Route{},
		echo:   e,
	}
}

func (r *Router) Add(method, path string, h HandlerFunc) {
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	pnames := []string{}
	ppath := path

	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1

			r.insert(method, path[:i], nil, skind, "", nil)
			for ; i < l && path[i] != '/'; i++ {

			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				r.insert(method, path[:i], h, pkind, ppath, pnames)
			} else {
				r.insert(method, path[:i], nil, pkind, "", nil)
			}
		} else if path[i] == '*' {
			r.insert(method, path[:i], nil, skind, "", nil)
			pnames = append(pnames, "*")
			r.insert(method, path[:i+1], h, akind, ppath, pnames)
		}
	}
	r.insert(method, path[:i+1], h, akind, ppath, pnames)
}

func (r *Router) insert(method, path string, h HandlerFunc, t kind, ppath string, pnames []string) {
	l := len(pnames)
	if *r.echo.maxParam < l {
		*r.echo.maxParam = l
	}

	cn := r.tree
	if cn == nil {
		panic("echo: invalid method")
	}
	search := path

	for {
		sl := len(search)
		pl := len(cn.prefix)
		l := 0

		max := pl
		if sl < max {
			max = sl
		}

		for ; l < max && search[l] == cn.prefix[l]; l++ {
		}

		if l == 0 {
			cn.label = search[0]
			cn.prefix = search
			if h != nil {
				cn.kind = t
				cn.addHandler(method, h)
				cn.ppath = ppath
				cn.pnames = pnames
			}
		} else if l < pl {
			n := newNode(cn.kind, cn.prefix[l:], cn, cn.children, cn.methodHandler, cn.ppath, cn.pnames)

			for _, child := range cn.children {
				child.parent = n
			}

			cn.kind = skind
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			cn.children = nil
			cn.methodHandler = new(methodHandler)
			cn.ppath = ""
			cn.pnames = nil

			cn.addChild(n)

			if l == sl {
				cn.kind = t
				cn.addHandler(method, h)
				cn.ppath = ppath
				cn.pnames = nil
			} else {
				n = newNode(t, search[l:], cn, nil, new(methodHandler), ppath, pnames)
				n.addHandler(method, h)
				cn.addChild(n)
			}
		} else if l < sl {
			search = search[l:]
			c := cn.findChildWithLabel(search[0])
			if c != nil {
				cn = c
				continue
			}
			n := newNode(t, search, cn, nil, new(methodHandler), ppath, pnames)
			n.addHandler(method, h)
			cn.addChild(n)
		} else {
			if h != nil {
				cn.addHandler(method, h)
				cn.ppath = ppath
				if len(cn.pnames) == 0 {
					cn.pnames = pnames
				}
			}
		}
		return
	}
}

func newNode(t kind, pre string, p *node, c children, mh *methodHandler, ppath string, pnames []string) *node {
	return &node{
		kind:          t,
		label:         pre[0],
		prefix:        pre,
		parent:        p,
		children:      c,
		ppath:         ppath,
		pnames:        pnames,
		methodHandler: mh,
	}
}

func (n *node) addChild(c *node) {
	n.children = append(n.children, c)
}

func (n *node) findChild(l byte, t kind) *node {
	for _, c := range n.children {
		if c.label == l && c.kind == t {
			return c
		}
	}
	return nil
}

func (n *node) findChildWithLabel(l byte) *node {
	for _, c := range n.children {
		if c.label == l {
			return c
		}
	}
	return nil
}

func (n *node) findChildByKind(t kind) *node {
	for _, c := range n.children {
		if c.kind == t {
			return c
		}
	}
	return nil
}

func (n *node) addHandler(method string, h HandlerFunc) {
	switch method {
	case http.MethodConnect:
		n.methodHandler.connect = h
	case http.MethodDelete:
		n.methodHandler.delete = h
	case http.MethodGet:
		n.methodHandler.get = h
	case http.MethodHead:
		n.methodHandler.head = h
	case http.MethodOptions:
		n.methodHandler.options = h
	case http.MethodPatch:
		n.methodHandler.patch = h
	case http.MethodPost:
		n.methodHandler.post = h
	case PROPFIND:
		n.methodHandler.propfind = h
	case http.MethodPut:
		n.methodHandler.put = h
	case http.MethodTrace:
		n.methodHandler.trace = h
	case REPORT:
		n.methodHandler.report = h
	}
}

func (n *node) finaHandler(method string) HandlerFunc {
	switch method {
	case http.MethodConnect:
		return n.methodHandler.connect
	case http.MethodDelete:
		return n.methodHandler.delete
	case http.MethodGet:
		return n.methodHandler.get
	case http.MethodHead:
		return n.methodHandler.head
	case http.MethodOptions:
		return n.methodHandler.options
	case http.MethodPatch:
		return n.methodHandler.patch
	case http.MethodPost:
		return n.methodHandler.post
	case PROPFIND:
		return n.methodHandler.propfind
	case http.MethodPut:
		return n.methodHandler.put
	case http.MethodTrace:
		return n.methodHandler.trace
	case REPORT:
		return n.methodHandler.report
	default:
		return nil
	}
}

func (n *node) checkMethodNotAllowed() HandlerFunc {
	for _, m := range methods {
		if h := n.finaHandler(m); h != nil {
			return MethodNotAllowedHandler
		}
	}
	return NotFoundHandler
}
