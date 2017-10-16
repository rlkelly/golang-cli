package main

import (
	"fmt"
	"os"
	"runtime"
  "strings"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

type sf func(int, string) string

func GetXY(x int, y int) (int, int) {
	if y == -1 {
		y = CurrentHeight() + 1
	}

	if x&PCT != 0 {
		x = int((x & 0xFF) * Width() / 100)
	}

	if y&PCT != 0 {
		y = int((y & 0xFF) * Height() / 100)
	}

	return x, y
}

func Height() int {
	ws, err := getWinsize()
	if err != nil {
		return -1
	}
	return int(ws.Row)
}

func Width() int {
	ws, err := getWinsize()

	if err != nil {
		return -1
	}

	return int(ws.Col)
}

func getWinsize() (*winsize, error) {
	ws := new(winsize)

	var _TIOCGWINSZ int64

	switch runtime.GOOS {
	case "linux":
		_TIOCGWINSZ = 0x5413
	case "darwin":
		_TIOCGWINSZ = 1074295912
	}

	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(_TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if int(r1) == -1 {
		fmt.Println("Error:", os.NewSyscallError("GetWinsize", errno))
		return nil, os.NewSyscallError("GetWinsize", errno)
	}
	return ws, nil
}

func CurrentHeight() int {
	return strings.Count(Screen.String(), "\n")
}

// func Print(a ...interface{}) (n int, err error) {
// 	return fmt.Fprint(Screen, a...)
// }

// func Println(a ...interface{}) (n int, err error) {
// 	return fmt.Fprintln(Screen, a...)
// }

func MoveTo(str string, x int, y int) (out string) {
	x, y = GetXY(x, y)

	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("\033[%d;%dH%s", y+idx, x, line)
	})
}

func applyTransform(str string, transform sf) (out string) {
	out = ""

	for idx, line := range strings.Split(str, "\n") {
		out += transform(idx, line)
	}

	return
}
