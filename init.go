//go:build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"
)

// https://github.com/sirupsen/logrus/issues/172#issuecomment-353724264

func init() {
	var originalMode uint32
	stdout := windows.Handle(os.Stdout.Fd())

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
