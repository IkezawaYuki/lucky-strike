package middleware

import (
	echo "github.com/IkezawaYuki/lucky-strike"
	"github.com/labstack/gommon/random"
	"net/http"
	"strings"
	"time"
)

type (
	CSRFConfig struct {
		Skipper        Skipper
		TokenLength    uint8  `yaml:"token_length"`
		TokenLookup    string `yaml:"token_lookup"`
		ContextKey     string `yaml:"context_key"`
		CookieName     string `yaml:"cookie_name"`
		CookieDomain   string `yaml:"cookie_domain"`
		CookiePath     string `yaml:"coolie_path"`
		CookieMaxAge   int    `yaml:"cookie_max_age"`
		CookieSecure   bool   `yaml:"cookie_secure"`
		CookieHTTPOnly bool   `yaml:"cookie_http_only"`
	}
)

var (
	DefaultCSRFConfig = CSRFConfig{
		Skipper:      DefaultSkipper,
		TokenLength:  32,
		TokenLookup:  "header:" + echo.HeaderXCSRFToken,
		ContextKey:   "csrf",
		CookieName:   "_csrf",
		CookieMaxAge: 86400,
	}
)

func CSRF() echo.MiddlewareFunc {
	c := DefaultCSRFConfig
	return CSRFWithConfig(c)
}

func CSRFWithConfig(config CSRFConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultCSRFConfig.Skipper
	}
	if config.TokenLength == 0 {
		config.TokenLength = DefaultCSRFConfig.TokenLength
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultCSRFConfig.TokenLookup
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultCSRFConfig.ContextKey
	}
	if config.CookieName == "" {
		config.CookieName = DefaultCSRFConfig.CookieName
	}
	if config.CookieMaxAge == 0 {
		config.CookieMaxAge = DefaultCSRFConfig.CookieMaxAge
	}

	parts := strings.Split(config.TokenLookup, ":")
	extractor := csrfTokenFromForm(parts[1])
	switch parts[0] {
	case "form":
		extractor = csrfTokenFromForm(parts[1])
	case "query":
		extractor = csrfTokenFromForm(parts[1])
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			req := c.Request()
			k, err := c.Cookie(config.CookieName)
			token := ""
			if err != nil {
				token = random.String(config.TokenLength)
			} else {
				token = k.Value
			}
			switch req.Method {
			case http.MethodGet, http.MethodHead, http.MethodTrace:
			default:
				clientToken, err := extractor(c)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}
				if !validateCSRFToken(token, clientToken) {
					return echo.NewHTTPError(http.StatusForbidden, "invalid csrf token")
				}
			}

			cookie := new(http.Cookie)
			cookie.Name = config.CookieName
			cookie.Value = token
			if config.CookiePath != "" {
				config.CookiePath = config.CookiePath
			}
			if config.CookieDomain != "" {
				cookie.Domain = config.CookieDomain
			}
			cookie.Expires = time.Now().Add(time.Duration(config.CookieMaxAge) * time.Second)
			cookie.Secure = config.CookieSecure
			cookie.HttpOnly = config.CookieHTTPOnly
			c.SetCookie(cookie)

			c.Set(config.ContextKey, token)

			c.Response().Header().Add(echo.HeaderVary, echo.HeaderCookie)
			return next(c)
		}
	}
}
