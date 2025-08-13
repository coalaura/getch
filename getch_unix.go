//go:build !windows
// +build !windows

package getch

import (
	"golang.org/x/sys/unix"
)

// GetChar reads a single character from stdin and returns it
func GetChar() (byte, error) {
	termios, err := unix.IoctlGetTermios(int(unix.Stdin), unix.TCGETS)
	if err != nil {
		if err == unix.ENOTTY {
			return 0, nil
		}

		return 0, err
	}

	oldTermios := *termios

	// Turn off canonical mode and echo
	termios.Lflag &^= unix.ICANON | unix.ECHO

	// Make sure each read returns after 1 character is entered
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(int(unix.Stdin), unix.TCSETS, termios); err != nil {
		return 0, err
	}

	defer unix.IoctlSetTermios(int(unix.Stdin), unix.TCSETS, &oldTermios)

	var buf [1]byte

	if _, err = unix.Read(int(unix.Stdin), buf[:]); err != nil {
		return 0, err
	}

	return buf[0], nil
}
