package middleware

import (
	"fmt"
	echo "github.com/IkezawaYuki/lucky-strike"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/copystructure"
	"net/http"
)

type (
	JWTConfig struct {
		Skipper        Skipper
		BeforeFunc     BeforeFunc
		SuccessHandler JWTSuccessHandler
		ErrorHandler   JWTErrorHandler
		SigningKey     interface{}
		SigningKeys    map[string]interface{}
		SigingMethod   string
		ContextKey     string
		Claims         jwt.Claims
		TokenLookup    string
		AuthScheme     string
		keyFunc        jwt.Keyfunc
	}
	JWTSuccessHandler         func(echo.Context)
	JWTErrorHandler           func(error) error
	JWTErrorHandleWithContext func(error, echo.Context) error
	jwtExtractor              func(echo.Context) (string, error)
)

const (
	AlgorithmHS256 = "HS256"
)

var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "missing or malformed jwt")
)

var (
	DefaultJWTConfig = JWTConfig{
		Skipper:      DefaultSkipper,
		SigingMethod: AlgorithmHS256,
		ContextKey:   "user",
		Claims:       jwt.MapClaims{},
		TokenLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme:   "Bearer",
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
		if t.Method.Alg() != config.SigingMethod {
			return nil, fmt.Errorf("unexpected jwt siging method=%v", t.Header["alg"])
		}

	}
}
