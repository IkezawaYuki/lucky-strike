package logging

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/panicwrap"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

const panicOutput = `
!!!!!!!!!!!!!!!!!!!!!!!!!!! TERRAFORM CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!
Terraform crashed! This is always indicative of a bug within Terraform.
A crash log has been placed at %[1]q relative to your current
working directory. It would be immensely helpful if you could please
report the crash with Terraform[1] so that we can fix this.
When reporting bugs, please include your terraform version. That
information is available on the first line of crash.log. You can also
get it by running 'terraform --version' on the command line.
SECURITY WARNING: the %[1]q file that was created may contain 
sensitive information that must be redacted before it is safe to share 
on the issue tracker.
[1]: https://github.com/hashicorp/terraform/issues
!!!!!!!!!!!!!!!!!!!!!!!!!!! TERRAFORM CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!
`

func PanicHandler(tmpLogPath string) panicwrap.HandlerFunc {
	return func(m string) {
		f, err := ioutil.TempFile(".", "crash.*.log")
		if err != nil {
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

		if _, err = io.Copy(f, tmpLog); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file %q: %v\n", tmpLogPath, err)
			return
		}

		f.WriteString("\n" + m)

		fmt.Printf("\n\n")
		fmt.Printf(panicOutput, f.Name())
	}
}

const pluginPanicOutput = `
Stack trace from the %[1]s plugin:
%s
Error: The %[1]s plugin crashed!
This is always indicative of a bug within the plugin. It would be immensely
helpful if you could report the crash with the plugin's maintainers so that it
can be fixed. The output above should help diagnose the issue.
`

func PluginPanics() []string {
	return panics.allPanics()
}

type panicRecorder struct {
	sync.Mutex

	panics map[string][]string

	maxLines int
}

func (p *panicRecorder) registerPlugin(name string) func(string) {
	p.Lock()
	defer p.Unlock()

	delete(p.panics, name)

	count := 0

	return func(line string) {
		p.Lock()
		defer p.Unlock()

		if count > p.maxLines {
			return
		}
		count++

		p.panics[name] = append(p.panics[name], line)
	}
}

func (p *panicRecorder) allPanics() []string {
	p.Lock()
	defer p.Unlock()

	var res []string
	for name, lines := range p.panics {
		if len(lines) == 0 {
			continue
		}
		res = append(res, fmt.Sprintf(pluginPanicOutput, name, strings.Join(lines, "\n")))
	}
	return res
}

type logPanicWrapper struct {
	hclog.Logger
	panicRecorder func(string)
	inPanic       bool
}

func (l *logPanicWrapper) Named(name string) hclog.Logger {
	return &logPanicWrapper{
		Logger:        l.Logger.Named(name),
		panicRecorder: panics.registerPlugin(name),
	}
}

func (l *logPanicWrapper) Debug(msg string, args ...interface{}) {
	panicPrefix := strings.HasPrefix(msg, "panic: ") || strings.HasPrefix(msg, "fatal error: ")

	l.inPanic = l.inPanic || panicPrefix

	if l.inPanic && l.panicRecorder != nil {
		l.panicRecorder(msg)
	}

	if panicPrefix {
		colon := strings.Index(msg, ":")
		msg = strings.ToUpper(msg[:colon]) + msg[colon:]
	}

	l.Logger.Debug(msg, args...)
}
