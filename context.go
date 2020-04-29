package echo

import (
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
)

type (
	Context interface {
		Request() *http.Request
		SetRequest(r *http.Request)
		SetResponse(r *Response)
		Response() *Response
		IsTLS() bool
		IsWebSocket() bool
		Scheme() string
		RealIP() string
		Path() string
		SetPath(p string)
		Param(name string) string
		ParamNames() []string
		ParamValues() []string
		SetParamValues(values ...string) []string
		QueryParam(name string) string
		QueryParams() url.Values
		QueryString() string
		FormValue(name string) string
		FormParams() (url.Values, error)
		FormFile(name string) (*multipart.FileHeader, error)
		MultipartForm() (*multipart.Form, error)
		Cookie(name string) (*http.Cookie, error)
		SetCookie(cookie *http.Cookie)
		Cookies() []*http.Cookie
		Get(key string) interface{}
		Set(key string, val interface{})
		Bind(i interface{}) error
		Validate(i interface{}) error
		Render(code int, name string, data interface{}) error
		HTML(code int, html string) error
		HTMLBlob(code int, b []byte) error
		String(code int, s string) error
		JSON(code int, i interface{}) error
		JSONPretty(code int, i interface{}, indent string) error
		JSONBlob(code int, b []byte) error
		JSONP(code int, callback string, i interface{}) error
		JSONPBlob(code int, callback string, b []byte) error
		XML(code int, i interface{}) error
		XMLPretty(code int, i interface{}, indent string) error
		XMLBlob(code int, b []byte) error
		Blob(code int, contentType string, b []byte) error
		Stream(code int, contentType string, r io.Reader) error
		File(file string) error
		Attachment(file string, name string) error
		Inline(file string, name string) error
		NoContent(code int) error
		Redirect(code int, url string) error
		Error(err error)
		Handler() HandlerFunc
		SetHandler(h HandlerFunc)
		Logger() Logger
		SetLogger(l Logger)
		Echo() *Echo
		Reset(r *http.Request, w http.ResponseWriter)
	}
	context struct {
		request  *http.Request
		response *Response
		path     string
		pnames   []string
		pvalues  []string
		query    url.Values
		handler  HandlerFunc
		store    Map
		echo     *Echo
		logger   Logger
		lock     sync.RWMutex
	}
)

const (
	defaultMemory = 32 << 20
	indexPage     = "index.html"
	defaultIndent = "  "
)

func (c *context) writeContentType(value string) {
	header := c.Response().Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) SetRequest(r *http.Request) {
	c.request = r
}

func (c *context) Response() *Response {
	return c.response
}

func (c *context) SetResponse(r *Response) {
	c.response = r
}

func (c *context) IsTLS() bool {
	return c.request.TLS != nil
}

func (c *context) IsWebSocket() bool {
	upgrade := c.request.Header.Get(HeaderUpgrade)
}

func (c *context) Logger() Logger {
	res := c.logger
	if res != nil {
		return res
	}
	return c.echo.Logger
}
