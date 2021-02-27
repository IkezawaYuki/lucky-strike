package main

import (
	"github.com/IkezawaYuki/lucky-strike/internal/logging"
	"github.com/mitchellh/panicwrap"
	"log"
	"os"
)

const (
	envTempLogPath = "TF_TEMP_LOG_PATH"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	var wrapConfig panicwrap.WrapConfig

	if os.Getenv("TF_FORK") == "0" {
		return wrappedMain()
	}
}

func wrappedMain() int {
	var err error

	tmpLogPath := os.Getenv(envTempLogPath)
	if tmpLogPath != "" {
		f, err := os.OpenFile(tmpLogPath, os.O_RDWR|os.O_APPEND, 0666)
		if err == nil {
			defer f.Close()

			log.Printf("[DEBUG] Adding temp file log sink: %s", f.Name())
			logging.RegisterSink(f)
		} else {
			// todo
		}
	}
}
