package echo

import (
	"fmt"
	"github.com/labstack/gommon/color"
	"golang.org/x/crypto/acme/autocert"
	"io"
	StdLog "log"
	"net"
	"net/http"
	"sync"
)

type (
	Echo struct {
		common
		StdLogger        StdLog.Logger
		colorer          *color.Color
		premiddleware    []MiddlewareFunc
		middleware       []MiddlewareFunc
		maxParam         *int
		router           *Router
		routers          map[string]*Router
		notFoundHandler  HandlerFunc
		pool             sync.Pool
		Server           *http.Server
		TLSServer        *http.Server
		Listener         net.Listener
		TLSListener      net.Listener
		AutoTLSManager   autocert.Manager
		DisableHTTP2     bool
		Debug            bool
		HideBanner       bool
		HidePort         bool
		HttpErrorHandler HTTPErrorHandler
		Binder           Binder
		Validator        Validator
		Renderer         Renderer
		Logger           Logger
		IPExtractor      IPExtractor
	}

	Route struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Name   string `json:"name"`
	}

	HTTPError struct {
		Code     int         `json:"-"`
		Message  interface{} `json:"message"`
		Internal error       `json:"-"`
	}
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	HandlerFunc func(Context) error

	HTTPErrorHandler func(error, Context)

	Validator interface {
		Validate(i interface{}) error
	}

	Renderer interface {
		Render(io.Writer, string, interface{}, Context) error
	}

	Map map[string]interface{}

	common struct{}
)

const (
	charsetUTF8 = "charset=UTF-8"
	PROPFIND    = "PROPFIND"
	REPORT      = "REPORT"
)

var (
	methods = [...]string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		PROPFIND,
		http.MethodPut,
		http.MethodTrace,
		REPORT,
	}
)

var (
	ErrUnsupportedMediaType = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrNotFound             = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized         = NewHTTPError(http.StatusUnauthorized)
	ErrForbidden            = NewHTTPError(http.StatusForbidden)
	ErrMethodNotAllowed     = NewHTTPError(http.StatusMethodNotAllowed)
)

var (
	NotFoundHandler = func(c Context) error {
		return ErrNotFound
	}
	MethodNotAllowedHandler = func(c Context) error {
		return ErrMethodNotAllowed
	}
)

func NewHTTPError(code int, message ...interface{}) *HTTPError {
	he := &HTTPError{
		Code:    code,
		Message: http.StatusText(code),
	}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}

func (he *HTTPError) Error() string {
	if he.Internal == nil {
		return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%v", he.Code, he.Message, he.Internal)
}

func (he *HTTPError) SetInterval(err error) *HTTPError {
	he.Internal = err
	return he
}
