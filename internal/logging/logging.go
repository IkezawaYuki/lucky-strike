package logging

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
)

const (
	envLog     = "TF_LOG"
	envLogFile = "TF_LOG_PATH"

	envLogCore = "TF_LOG_CORE"
)

var (
	ValidLevels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "OFF"}
	logger      hclog.Logger
	logWriter   io.Writer
	panics      = &panicRecorder{
		panic:    make(map[string][]string),
		maxLines: 100,
	}
)

func init() {
	logger = newHCLogger("")
	logWriter = logger.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(logWriter)
}

func RegisterSink(f *os.File) {
	l, ok := logger.(hclog.InterceptLogger)
	if !ok {
		panic("global logger is not an InterceptLogger")
	}

	if f == nil {
		return
	}

	l.RegisterSink(hclog.NewSinkAdapter(&hclog.LoggerOptions{
		Level:  hclog.Trace,
		Output: f,
	}))
}

func newHCLogger(name string) hclog.Logger {
	logOutput := io.Writer(os.Stderr)
	logLevel := globalLogLevel()

	if logPath := os.Getenv(envLogFile); logPath != "" {
		f, err := os.OpenFile(logPath, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		} else {
			logOutput = f
		}
	}

	return hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:              name,
		Level:             logLevel,
		Output:            logOutput,
		IndependentLevels: true,
	})
}

func globalLogLevel() hclog.Level {
	envLevel := strings.ToUpper(os.Getenv(envLog))
	if envLevel == "" {
		envLevel = strings.ToUpper(os.Getenv(envLogCore))
	}
	return parseLogLevel(envLevel)
}

func parseLogLevel(envLevel string) hclog.Level {
	if envLevel == "" {
		return hclog.Off
	}

	logLevel := hclog.Trace
	if isValidLogLevel(envLevel) {
		logLevel = hclog.LevelFromString(envLevel)
	} else {
		fmt.Fprintf(os.Stderr, "[WARN] Invalid log level: %q. Defaulting to level: TRACE. Valid levels are: %+v",
			envLevel, ValidLevels)
	}
	return logLevel
}

func isValidLogLevel(level string) bool {
	for _, l := range ValidLevels {
		if level == string(l) {
			return true
		}
	}
	return false
}
