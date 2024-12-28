//go:build windows
// +build windows

package getch

import (
	"golang.org/x/sys/windows"
)

// GetChar reads a single character from stdin and returns it
func GetChar() (byte, error) {
	handle, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil {
		if err == windows.ERROR_INVALID_HANDLE {
			return 0, nil
		}

		return 0, err
	}

	var mode uint32

	err = windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return 0, err
	}

	oldMode := mode

	// Disable line input and echo
	mode &^= windows.ENABLE_LINE_INPUT
	mode &^= windows.ENABLE_ECHO_INPUT
	mode &^= windows.ENABLE_PROCESSED_INPUT

	err = windows.SetConsoleMode(handle, mode)
	if err != nil {
		return 0, err
	}

	defer windows.SetConsoleMode(handle, oldMode)

	var (
		buf  [1]byte
		read uint32
	)

	if err = windows.ReadFile(handle, buf[:], &read, nil); err != nil {
		return 0, err
	}

	return buf[0], nil
}
