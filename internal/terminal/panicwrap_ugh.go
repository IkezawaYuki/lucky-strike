package terminal

func ReinitInsidePanicwrap(state *PrePanicwrapState) (*Streams, error) {
	ret, err := Init()
}

type PrePanicwrapState struct {
	StderrIsTerminal bool
	StderrWidth      int
}
