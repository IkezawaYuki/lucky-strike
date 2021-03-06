package terminal

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"golang.org/x/sys/windows"
	"os"
)

func configureOutputHandle(f *os.File) (*OutputStream, error) {
	ret := &OutputStream{
		File: f,
	}
	if fd := f.Fd(); isatty.IsTerminal(fd) {
		err := SetConsoleOutputCP(CP_UTF8)
		if err != nil {
			return nil, fmt.Errorf("failed to set the console to UTF-8 mode; you may need to use a newer version of Windows: %s", err)
		}

		ret.getColumns = getColumnsWindowsConsole
		var mode uint32
		err = windows.SetConsoleMode(windows.Handle(fd), mode)
		if err != nil {
			return ret, nil
		}

	}
}
