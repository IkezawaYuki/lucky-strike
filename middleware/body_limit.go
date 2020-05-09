package middleware

import (
	"fmt"
	echo "github.com/IkezawaYuki/lucky-strike"
	"github.com/labstack/gommon/bytes"
	"io"
	"sync"
)

type (
	BodyLimitConfig struct {
		Skipper Skipper
		Limit   string `yaml:"limit"`
		limit   int64
	}
	limitedReader struct {
		BodyLimitConfig
		reader  io.ReadCloser
		read    int64
		context echo.Context
	}
)

var (
	DefaultBodyLimitConfig = BodyLimitConfig{
		Skipper: DefaultSkipper,
	}
)

func BodyLimit(limit string) echo.MiddlewareFunc {
	c := DefaultBodyLimitConfig
	c.Limit = limit
	return BodyLimitWithConfig(c)
}

func BodyLimitWithConfig(config BodyLimitConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultBodyLimitConfig.Skipper
	}
	limit, err := bytes.Parse(config.Limit)
	if err != nil {
		panic(fmt.Errorf("echo: invalid body-limit=%s", config.Limit))
	}
	config.limit = limit
	pool := limitedReaderPool(config)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			req := c.Request()

			if req.ContentLength > config.limit {
				return echo.ErrStatusRequestEntityTooLarge
			}

			r := pool.Get().(*limitedReader)
			r.Reset(req.Body, c)
			defer pool.Put(r)
			req.Body = r
			return next(c)
		}
	}
}

func (r *limitedReader) Read(b []byte) (n int, err error) {
	n, err = r.reader.Read(b)
	r.read += int64(n)
	if r.read > r.limit {
		return n, echo.ErrStatusRequestEntityTooLarge
	}
	return
}

func (r *limitedReader) Close() error {
	return r.reader.Close()
}

func (r *limitedReader) Reset(reader io.ReadCloser, context echo.Context) {
	r.reader = reader
	r.context = context
	r.read = 0
}

func limitedReaderPool(c BodyLimitConfig) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return &limitedReader{
				BodyLimitConfig: c}
		},
	}
}
