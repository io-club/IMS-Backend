package util

import (
	"path"
	"runtime"
)

// Get the calling function's directory path
func GetCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1) //1 means the current function, 2 means the function that calls this function, and so on
	return path.Dir(filename)
}

// Get the current running code's file name and line number
func GetFileAndLineNo(skip int) (string, string, int) {
	funcName, file, line, ok := runtime.Caller(skip)
	if ok {
		return file, runtime.FuncForPC(funcName).Name(), line
	}
	return "", "", 0
}
