package middleware

import (
	"encoding/base64"
	echo "github.com/IkezawaYuki/lucky-strike"
	"strconv"
	"strings"
)

type (
	BasicAuthConfig struct {
		Skipper   Skipper
		Validator BasicAuthValidator
		Realm     string
	}
	BasicAuthValidator func(string, string, echo.Context) (bool, error)
)

const (
	basic        = "basic"
	defaultRealm = "Restricted"
)

var (
	DefaultBasicAuthConfig = BasicAuthConfig{
		Skipper: DefaultSkipper,
		Realm:   defaultRealm,
	}
)

func BasicAuth(fn BasicAuthValidator) echo.MiddlewareFunc {
	c := DefaultBasicAuthConfig
	c.Validator = fn
	return BasicAuthWithConfig(c)
}

func BasicAuthWithConfig(config BasicAuthConfig) echo.MiddlewareFunc {
	if config.Validator == nil {
		panic("echo: basic-auth middleware requires a validator function")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultBasicAuthConfig.Skipper
	}
	if config.Realm == "" {
		config.Realm = defaultRealm
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			l := len(basic)

			if len(auth) > l+1 && strings.ToLower(auth[:l]) == basic {
				b, err := base64.StdEncoding.DecodeString(auth[l+1:])
				if err != nil {
					return err
				}
				cred := string(b)
				for i := 0; i < len(cred); i++ {
					if cred[i] == ':' {
						valid, err := config.Validator(cred[:i], cred[i+1:], c)
						if err != nil {
							return err
						} else if valid {
							return next(c)
						}
						break
					}
				}
			}

			realm := defaultRealm
			if config.Realm != defaultRealm {
				realm = strconv.Quote(config.Realm)
			}

			c.Response().Header().Set(echo.HeaderWWWAuthenticate, basic+" realm="+realm)
			return echo.ErrUnauthorized
		}
	}
}
