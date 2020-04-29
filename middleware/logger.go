package middleware

import (
	echo "github.com/IkezawaYuki/lucky-strike"
	"github.com/labstack/gommon/color"
	"github.com/valyala/fasttemplate"
	"io"
	"sync"
)

type (
	LoggerConfig struct {
		Skipper          Skipper
		Format           string `yaml:"format"`
		CustomTimeFormat string `yaml:"custom_time_format"`
		Output           io.Writer
		template         *fasttemplate.Template
		colorer          *color.Color
		pool             *sync.Pool
	}
)

var (
	DefaultLoggerConfig = LoggerConfig{
		Skipper: DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		colorer:          color.New(),
	}
)

func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {

}
