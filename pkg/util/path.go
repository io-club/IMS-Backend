package util

import (
	"path"
	"runtime"
)

// 获取调用方的的地址路径
func GetCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1) //1 means the current function, 2 means the function that calls this function, and so on
	return path.Dir(filename)
}

// 获取当前运行代码的文件名和行号
func GetFileAndLineNo(skip int) (string, string, int) {
	funcName, file, line, ok := runtime.Caller(skip)
	if ok {
		return file, runtime.FuncForPC(funcName).Name(), line
	}
	return "", "", 0
}
