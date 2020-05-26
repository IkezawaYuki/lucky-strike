package middleware

import (
	echo "github.com/IkezawaYuki/lucky-strike"
	"net/http"
)

type (
	MethodOverrideConfig struct {
		Skipper Skipper
		Getter  MethodOverrideGetter
	}
	MethodOverrideGetter func(echo.Context) string
)

var (
	DefaultMethodOverrideConfig = MethodOverrideConfig{
		Skipper: DefaultSkipper,
		Getter:  MethodFromHeader(echo.HeaderXHTTPMethodOverride),
	}
)

func MethodOverride() echo.MiddlewareFunc {
	return MethodOverrideWithConfig(DefaultMethodOverrideConfig)
}

func MethodOverrideWithConfig(config MethodOverrideConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultMethodOverrideConfig.Skipper
	}
	if config.Getter == nil {
		config.Getter = DefaultMethodOverrideConfig.Getter
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			req := c.Request()
			if req.Method == http.MethodPost {
				m := config.Getter(c)
				if m != "" {
					req.Method = m
				}
			}
			return next(c)
		}
	}
}

func MethodFromHeader(header string) MethodOverrideGetter {
	return func(c echo.Context) string {
		return c.Request().Header.Get(header)
	}
}

func MethodFromForm(param string) MethodOverrideGetter {
	return func(c echo.Context) string {
		return c.FormValue(param)
	}
}

func MethodFromQuery(param string) MethodOverrideGetter {
	return func(c echo.Context) string {
		return c.QueryParam(param)
	}
}
