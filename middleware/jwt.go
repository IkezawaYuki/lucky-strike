package middleware

import (
	echo "github.com/IkezawaYuki/lucky-strike"
	"github.com/dgrijalva/jwt-go"
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
