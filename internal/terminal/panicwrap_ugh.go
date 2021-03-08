package terminal

import (
	"os"
)

func ReinitInsidePanicwrap(state *PrePanicwrapState) (*Streams, error) {
	ret, err := Init()
	if err != nil {
		return ret, err
	}
	if state != nil {
		ret.Stderr = &OutputStream{
			File: ret.Stderr.File,
			isTerminal: func(f *os.File) bool {
				return state.StderrIsTerminal
			},
			getColumns: func(f *os.File) int {
				return state.StderrWidth
			},
		}
	}
	return ret, nil
}

type PrePanicwrapState struct {
	StderrIsTerminal bool
	StderrWidth      int
}
