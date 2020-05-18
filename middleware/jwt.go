package middleware

import (
	"fmt"
	echo "github.com/IkezawaYuki/lucky-strike"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/copystructure"
	"net/http"
	"reflect"
	"strings"
)

type (
	JWTConfig struct {
		Skipper                 Skipper
		BeforeFunc              BeforeFunc
		SuccessHandler          JWTSuccessHandler
		ErrorHandler            JWTErrorHandler
		ErrorHandlerWithContext JWTErrorHandlerWithContext
		SigningKey              interface{}
		SigningKeys             map[string]interface{}
		SigningMethod           string
		ContextKey              string
		Claims                  jwt.Claims
		TokenLookup             string
		AuthScheme              string
		keyFunc                 jwt.Keyfunc
	}
	JWTSuccessHandler          func(echo.Context)
	JWTErrorHandler            func(error) error
	JWTErrorHandlerWithContext func(error, echo.Context) error
	jwtExtractor               func(echo.Context) (string, error)
)

const (
	AlgorithmHS256 = "HS256"
)

var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "missing or malformed jwt")
)

var (
	DefaultJWTConfig = JWTConfig{
		Skipper:       DefaultSkipper,
		SigningMethod: AlgorithmHS256,
		ContextKey:    "user",
		Claims:        jwt.MapClaims{},
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
	}
)

func JWT(key interface{}) echo.MiddlewareFunc {
	c := DefaultJWTConfig
	c.SigningKey = key
	return JWTWithConfig(c)
}

func JWTWithConfig(config JWTConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultJWTConfig.Skipper
	}
	if config.SigningKey == nil && len(config.SigningKeys) == 0 {
		panic("echo: jwt middleware requires signing key")
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultJWTConfig.ContextKey
	}
	if config.Claims == nil {
		config.Claims = DefaultJWTConfig.Claims
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultJWTConfig.TokenLookup
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultJWTConfig.TokenLookup
	}
	config.keyFunc = func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != config.SigningKey {
			return nil, fmt.Errorf("unexpected jwt siging method=%v", t.Header["alg"])
		}
		if len(config.SigningKeys) > 0 {
			if kid, ok := t.Header["kid"].(string); ok {
				if key, ok := config.SigningKeys[kid]; ok {
					return key, nil
				}
			}
			return nil, fmt.Errorf("unexpected jwt key id=%v", t.Header["kid"])
		}
		return config.SigningKey, nil
	}

	parts := strings.Split(config.TokenLookup, ":")
	extractor := jwtFromHeader(parts[1], config.AuthScheme)
	switch parts[0] {
	case "query":
		extractor = jwtFromQuery(parts[1])
	case "param":
		extractor = jwtFromParam(parts[1])
	case "cookie":
		extractor = jwtFromCookie(parts[1])
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			if config.BeforeFunc != nil {
				config.BeforeFunc(c)
			}

			auth, err := extractor(c)
			if err != nil {
				if config.ErrorHandler != nil {
					return config.ErrorHandler(err)
				}
				if config.ErrorHandlerWithContext != nil {
					return config.ErrorHandlerWithContext(err, c)
				}
				return err
			}
			token := new(jwt.Token)
			if _, ok := config.Claims.(jwt.MapClaims); ok {
				token, err = jwt.Parse(auth, config.keyFunc)
			} else {
				t := reflect.ValueOf(config.Claims).Type().Elem()
				claims := reflect.New(t).Interface().(jwt.Claims)
				token, err = jwt.ParseWithClaims(auth, claims, config.keyFunc)
			}
			if err == nil && token.Valid {
				c.Set(config.ContextKey, token)
				if config.SuccessHandler != nil {
					config.SuccessHandler(c)
				}
				return next(c)
			}
			if config.ErrorHandler != nil {
				return config.ErrorHandler(err)
			}
			if config.ErrorHandlerWithContext != nil {
				return config.ErrorHandlerWithContext(err, c)
			}
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "invalid or expired jwt",
				Internal: err,
			}
		}
	}
}

func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}

func jwtFromQuery(param string) jwtExtractor {
	return func(c echo.Context) (string, error) {

	}
}
