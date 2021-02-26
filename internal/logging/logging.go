package logging

import (
	"github.com/hashicorp/go-hclog"
	"io"
	"os"
	"strings"
)

const (
	envLog = "TF_LOG"

	envLogCore = "TF_LOG_PATH"
)

func init() {
	logger = newHCLogger("")
}

func newHCLogger(name string) hclog.Logger {
	logOutput := io.Writer(os.Stderr)
	logLegel := globalLogLevel()
}

func globalLogLevel() hclog.Level {
	envLevel := strings.ToUpper(os.Getenv(envLog))
	if envLevel == "" {
		envLevel = strings.ToUpper(os.Getenv(envLogCore))
	}
	return parseLogLevel(envLevel)
}

func parseLogLevel(envLevel string) hclog.Level {
	if
}
