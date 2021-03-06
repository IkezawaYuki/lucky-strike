package terminal

import "os"

func configureOutputHandle(f *os.File) (*OutputStream, error) {
	return &OutputStream{
		File:       f,
		isTerminal: isTerminalGolangXTerm,
		getColumns: getColumnsGolangXTerm,
	}, nil
}

func configureInputHandle(f *os.File) (*InputStream, error) {
	return &InputStream{
		File:       f,
		isTerminal: isTerminalGolangXTerm,
	}, nil
}

func isTerminalGolangXTerm(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

func getColumnsGolangXTerm(f *os.File) int {
	width, _, err := term.GetSize(int(f.Fd()))
	if err != nil {
		return defaultColumns
	}
	return width
}
