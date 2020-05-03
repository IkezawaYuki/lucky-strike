package echo

import (
	"bufio"
	"net"
	"net/http"
)

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
		r.echo.Logger.Warn("response already committed")
	}
	for _, fn := range r.beforeFuncs {
		fn()
	}
	r.Status = code
	r.Writer.WriteHeader(code)
	r.Committed = true
}

func (r *Response) Write(b []byte) (n int, err error) {
	if !r.Committed {
		if r.Status == 0 {
			r.Status = http.StatusOK
		}
		r.WriteHeader(r.Status)
	}
	n, err = r.Writer.Write(b)
	r.Size += int64(n)
	for _, fn := range r.afterFuncs {
		fn()
	}
	return
}

func (r *Response) Flush() {
	r.Writer.(http.Flusher).Flush()
}

func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.Writer.(http.Hijacker).Hijack()
}

func (r *Response) reset(w http.ResponseWriter) {
	r.beforeFuncs = nil
	r.afterFuncs = nil
	r.Writer = w
	r.Size = 0
	r.Status = http.StatusOK
	r.Committed = false
}
