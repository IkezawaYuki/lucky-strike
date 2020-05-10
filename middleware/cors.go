package middleware

type (
	CORSConfig struct {
		Skipper       Skipper
		AllowOrigins  []string `yaml:"allow_origins"`
		AllowMethods  []string `yaml:"allow_methods"`
		AllowHeaders  []string `yaml:"allow_headers"`
		ExposeHeaders []string `yaml:"expose_headers"`
		MaxAge        int      `yaml:"max_age"`
	}
)
