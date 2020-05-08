package middleware

import (
	echo "github.com/IkezawaYuki/lucky-strike"
	"io"
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
