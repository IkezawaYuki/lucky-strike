package logging

import (
	"fmt"
	"github.com/mitchellh/panicwrap"
	"io"
	"io/ioutil"
	"os"
)

func PanicHandler(tmpLogPath string) panicwrap.HandlerFunc {
	return func(m string) {
		f, err := ioutil.TempFile(".", "crash.*.log")
		if err != nil{
			fmt.Fprintf(os.Stderr, "Failed to create crash log file: %s", err)
			return
		}
		defer f.Close()

		tmpLog, err := os.Open(tmpLogPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file %q: %v\n", tmpLogPath, err)
			return
		}
		defer tmpLog.Close()

		if _, err = io.Copy(f, tmpLog);
	}
}
