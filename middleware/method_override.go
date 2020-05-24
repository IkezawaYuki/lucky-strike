package middleware

import echo "github.com/IkezawaYuki/lucky-strike"

type (
	MethodOverrideConfig struct {
		Skipper Skipper
		Getter  MethodOverrideConfig
	}
	MethodOverrideGetter func(echo.Context) string
)

var (
	DefaultMethodOverriddGettter func(echo.Context) string
)
