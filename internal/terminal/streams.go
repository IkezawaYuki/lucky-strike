package terminal

import (
	"fmt"
	"os"
)

type Streams struct {
	Stdout *OutputStream
	Stderr *OutputStream
	Stdin  *InputStream
}

func Init() (*Streams, error) {
	stderr, err := configureOutputHandle(os.Stderr)
	if err != nil {
		return nil, err
	}

	stdout, err := configureOutputHandle(os.Stdout)
	if err != nil {
		return nil, err
	}

	stdin, err := configureInputHandle(os.Stdin)
	if err != nil {
		return nil, err
	}

	return &Streams{
		Stdout: stdout,
		Stderr: stderr,
		Stdin:  stdin,
	}, nil
}

func (s *Streams) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(s.Stdout.File, a...)
}

func (s *Streams) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(s.Stdout.File, format, a...)
}

func (s *Streams) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(s.Stdout.File, a...)
}

func (s *Streams) Eprint(a ...interface{}) (n int, err error) {
	return fmt.Fprint(s.Stderr.File, a...)
}

func (s *Streams) Eprintf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(s.Stderr.File, format, a...)
}

func (s *Streams) Eprintln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(s.Stderr.File, a...)
}
