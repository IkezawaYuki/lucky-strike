package terminal

import "os"

func configureOutputHandle(f *os.File) (*OutputStream, error) {
	return &OutputStream{
		File:       f,
		isTerminal: isTerminalGoralngXTerm,
		getColumns: getColumnsGolangXTerm,
	}, nil
}

func configureInputHandle(f *os.File) (*InputStream, error) {
	return &InputStream{
		File:       f,
		isTerminal: isTerminalGolangXTerm,
	}, nil
}
