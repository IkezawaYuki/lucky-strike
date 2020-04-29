package echo

import "net/http"

type (
	Response struct {
		echo        *Echo
		beforeFuncs []func()
		afterFuncs  []func()
		Writer      http.ResponseWriter
		Status      int
		Size        int64
		Committed   bool
	}
)

func NewResponse(w http.ResponseWriter, e *Echo) (r *Response) {
	return &Response{
		Writer: w,
		echo:   e,
	}
}

func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

func (r *Response) Before(fn func()) {
	r.beforeFuncs = append(r.beforeFuncs, fn)
}

func (r *Response) After(fn func()) {
	r.afterFuncs = append(r.afterFuncs, fn)
}

func (r *Response) WriteHeader(code int) {
	if r.Committed {
		r.echo.Logger
	}
}
