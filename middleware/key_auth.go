package middleware

import (
	"errors"
	echo "github.com/IkezawaYuki/lucky-strike"
	"net/http"
	"strings"
)

type (
	KeyAuthConfig struct {
		Skipper    Skipper
		KeyLookup  string `yaml:"key_lookup"`
		AuthScheme string
		Validator  KeyAuthValidator
	}
	KeyAuthValidator func(string, echo.Context) (bool, error)
	keyExtractor     func(echo.Context) (string, error)
)

var (
	DefaultKeyAuthConfig = KeyAuthConfig{
		Skipper:    DefaultSkipper,
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
	}
)

func KeyAuthWithConfig(config KeyAuthConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultKeyAuthConfig.Skipper
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultKeyAuthConfig.AuthScheme
	}
	if config.Validator == nil {
		panic("echo: key-auth middleware requires a validator function")
	}
	parts := strings.Split(config.KeyLookup, ":")
	extracor := keyFromHeader(parts[1], config.AuthScheme)
	switch parts[0] {
	case "query":
		extracor = keyFromQuery(parts[1])
	case "form":
		extracor = keyFromForm(parts[1])
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			key, err := extracor(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			valid, err := config.Validator(key, c)
			if err != nil {
				return &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  "invalid key",
					Internal: err,
				}
			} else if valid {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}

func keyFromHeader(header string, authScheme string) keyExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		if auth == "" {
			return "", errors.New("missing key in request header")
		}
		if header == echo.HeaderAuthorization {
			l := len(authScheme)
			if len(auth) > l+1 && auth[:l] == authScheme {
				return auth[l+1:], nil
			}
			return "", errors.New("invalid key in the request header")
		}
		return auth, nil
	}
}

func keyFromQuery(param string) keyExtractor {
	return func(c echo.Context) (string, error) {
		key := c.QueryParam(param)
		if key == "" {
			return "", errors.New("missing key in the query string")
		}
		return key, nil
	}
}

func keyFromForm(param string) keyExtractor {
	return func(c echo.Context) (string, error) {
		key := c.FormValue(param)
		if key == "" {
			return "", errors.New("missing key in the form")
		}
		return key, nil
	}
}
