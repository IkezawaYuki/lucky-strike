package terminal

import "os"

func configureOutputHandle(f *os.File) (*OutputStream, error) {
	return &OutputStream{
		File:       f,
		isTerminal: isTerminalGoralngXTerm,
		getColumns: getColumnsGolangXTerm,
	}, nil
}
