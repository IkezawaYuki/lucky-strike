package terminal

type Streams struct {
	Stdout *OutputStream
	Stderr *OutputStream
	Stdin  *InputStream
}
