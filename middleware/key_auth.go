package middleware

import echo "github.com/IkezawaYuki/lucky-strike"

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
