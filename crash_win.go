// +build windows

package log

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func setStdHandle(stdhandle int32, handle syscall.Handle) error {
	r0, _, e1 := syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(stdhandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}

func CrashLog(file string) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err.Error())
	} else {
		err = setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(f.Fd()))
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func SetConsoleColor(t LogType) string {
	mod := syscall.NewLazyDLL("kernel32.dll")
	SetConsoleTextAttribute := mod.NewProc("SetConsoleTextAttribute")
	getStdHandler, _ := syscall.GetStdHandle(-11)

	logStr, logColor := LogTypeToString(t)
	colorAttribute, err := strconv.ParseInt(strings.Replace(logColor, "0x", "", -1), 16, 64)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// set console color
	SetConsoleTextAttribute.Call(uintptr(getStdHandler), uintptr(colorAttribute))
	return logStr
}

func ResetConsoleColor() {
	mod := syscall.NewLazyDLL("kernel32.dll")
	SetConsoleTextAttribute := mod.NewProc("SetConsoleTextAttribute")
	getStdHandler, _ := syscall.GetStdHandle(-11)
	SetConsoleTextAttribute.Call(uintptr(getStdHandler), uintptr(0x07))
}
