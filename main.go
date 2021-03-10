package main

import (
	"fmt"
	"github.com/IkezawaYuki/lucky-strike/internal/logging"
	"github.com/IkezawaYuki/lucky-strike/internal/terminal"
	"github.com/hashicorp/terraform/version"
	"github.com/mitchellh/panicwrap"
	"log"
	"os"
	"runtime"
)

var Version = version.Version

var VersionPrerelease = version.Prerelease

const (
	envTempLogPath                 = "TF_TEMP_LOG_PATH"
	envTerminalPanicwrapWorkaround = "TF_PANICWRAP_STDERR"
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
			log.Printf("[ERROR] Could not open tmp log file: %v", err)
		}
	}

	log.Printf("[INFO] Terraform version: %s %s",
		Version, VersionPrerelease)
	log.Printf("[INFO] Go runtime version: %s", runtime.Version())
	log.Printf("[INFO] CLI args: %#v", os.Args)

	var streamState *terminal.PrePanicwrapState
	if raw := os.Getenv(envTerminalPanicwrapWorkaround); raw != "" {
		streamState = &terminal.PrePanicwrapState{}
		if _, err := fmt.Sscan(raw, "%t:%d", &streamState.StderrIsTerminal, &streamState.StderrWidth); err != nil {
			log.Printf("[WARN] %s is set but is incorrectly-formatted: %s", envTerminalPanicwrapWorkaround, err)
			streamState = nil
		}
		streams, err := terminal.ReinitInsidePanicwrap(streamState)
		if err != nil {
			Ui.Error
		}
	}
}
