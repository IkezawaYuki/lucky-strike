package middleware

import (
	"bytes"
	echo "github.com/IkezawaYuki/lucky-strike"
	"io"
	"io/ioutil"
	"net/http"
)

type (
	BodyDumpConfig struct {
		Skipper Skipper
		Handler BodyDumpHandler
	}
	BodyDumpHandler func(echo.Context, []byte, []byte)

	bodyDumpResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

var (
	DefaultBodyDumpConfig = BodyDumpConfig{
		Skipper: DefaultSkipper,
	}
)

func BodyDump(handler BodyDumpHandler) echo.MiddlewareFunc {
	c := DefaultBodyDumpConfig
	c.Handler = handler
	return BodyDumpWithConfig(c)
}

func BodyDumpWithConfig(config BodyDumpConfig) echo.MiddlewareFunc {
	if config.Handler == nil {
		panic("echo: body-dump middleware requires a handler function")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultBodyDumpConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			reqBody := []byte{}
			if c.Request().Body != nil {
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{
				Writer:         mw,
				ResponseWriter: c.Response().Writer,
			}
			c.Response().Writer = writer
			if err = next(c); err != nil {
				c.Error(err)
			}
			config.Handler(c, reqBody, resBody.Bytes())
			return
		}
	}
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
