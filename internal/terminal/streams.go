package terminal

import "os"

type Streams struct {
	Stdout *OutputStream
	Stderr *OutputStream
	Stdin  *InputStream
}

func Init() (*Streams, error) {
	stderr, err := configureOutputHandle(os.Stderr)
}
