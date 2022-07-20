package logger

import (
	"runtime"
	"strings"
)

func getLocation() (function, file string, line int) {
	pc, file, line, _ := runtime.Caller(3)
	function = getFuncName(pc)
	return
}

func getFuncName(pc uintptr) string {
	fn := runtime.FuncForPC(pc).Name()
	splitted := strings.Split(fn, "/")
	return splitted[len(splitted)-1]
}
