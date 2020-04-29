package echo

import (
	"mime/multipart"
	"net/http"
	"net/url"
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

		handler HandlerFunc
	}
)
