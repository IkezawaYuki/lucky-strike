package middleware

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
