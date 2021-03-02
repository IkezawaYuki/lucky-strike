package terminal

import "os"

const defaultColumns int = 78
const defaultIsTerminal bool = false

type OutputStream struct {
	File       *os.File
	isTerminal func(file *os.File) bool
	getColumns func(file *os.File) int
}

func (s *OutputStream) Columns() int {
	if s.getColumns == nil {
		return defaultColumns
	}
}
