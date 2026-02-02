// Based on chi's terminal helpers.
// https://github.com/go-chi/chi

package middleware

import (
	"fmt"
	"io"
	"os"
)

// IsTTY reports whether stdout is connected to a terminal.
var IsTTY bool

func init() {
	fi, err := os.Stdout.Stat()
	if err == nil {
		m := os.ModeDevice | os.ModeCharDevice
		IsTTY = fi.Mode()&m == m
	}
}

// ANSI color codes — normal.
var (
	nBlack   = []byte{'\033', '[', '3', '0', 'm'}
	nYellow  = []byte{'\033', '[', '3', '3', 'm'}
)

// ANSI color codes — bright/bold.
var (
	bRed     = []byte{'\033', '[', '3', '1', ';', '1', 'm'}
	bGreen   = []byte{'\033', '[', '3', '2', ';', '1', 'm'}
	bBlue    = []byte{'\033', '[', '3', '4', ';', '1', 'm'}
	bMagenta = []byte{'\033', '[', '3', '5', ';', '1', 'm'}
	bCyan    = []byte{'\033', '[', '3', '6', ';', '1', 'm'}
	bWhite   = []byte{'\033', '[', '3', '7', ';', '1', 'm'}

	reset = []byte{'\033', '[', '0', 'm'}
)

// cW writes a coloured string to w when IsTTY and useColor are true.
func cW(w io.Writer, useColor bool, color []byte, s string, args ...interface{}) {
	if IsTTY && useColor {
		w.Write(color)
	}
	fmt.Fprintf(w, s, args...)
	if IsTTY && useColor {
		w.Write(reset)
	}
}
