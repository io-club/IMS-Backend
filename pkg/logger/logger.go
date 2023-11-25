package iologger

import (
	"fmt"
	ioconfig "ims-server/pkg/config"
	ioconst "ims-server/pkg/consts"
	"ims-server/pkg/util"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger

	logLevel = 0 // 默认的 LogLevel 为 0，即所有级别的日志都打印
)

var (
	logOut    *os.File      // 日志输出文件
	heartbeat time.Duration // 日志心跳间隔
	maxAge    time.Duration // 日志保留时间
	logFile   string        // 日志文件名称

	timeFlag      time.Time    // 标记上次日志切换时间
	dayChangeLock sync.RWMutex // 切换日志文件时上锁
)

func init() {
	var ok bool
	var err error

	logLevel, ok = ioconst.LoggerLevel_value[ioconfig.GetLoggerConf().Level]
	if !ok {
		panic("获取配置\"日志级别\"发生错误")
	}

	logFile = ioconfig.GetLoggerConf().FileName
	postFix := "_" + time.Now().Format("2006-01-02 15-04-05")
	file := filepath.Join(ioconfig.GetLoggerConf().Path, logFile+postFix)
	logOut, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
	debugLogger.SetOutput(io.MultiWriter(debugLogger.Writer(), os.Stdout)) // 将输出同步到终端

	infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
	infoLogger.SetOutput(io.MultiWriter(infoLogger.Writer(), os.Stdout))

	warnLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
	warnLogger.SetOutput(io.MultiWriter(warnLogger.Writer(), os.Stdout))

	errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
	errorLogger.SetOutput(io.MultiWriter(errorLogger.Writer(), os.Stdout))

	heartbeat = ioconfig.GetLoggerConf().HeartBeat
	maxAge = ioconfig.GetLoggerConf().MaxAge

	timeFlag = time.Now()
	dayChangeLock = sync.RWMutex{}
}

func checkLogRotation() {
	now := time.Now()
	// 大部分时候会直接 return，所以并不慢
	if now.Sub(timeFlag) <= heartbeat {
		return
	}
	// 需要切换日志文件，进行上锁
	dayChangeLock.Lock()
	defer dayChangeLock.Unlock()

	err := logOut.Close()
	if err != nil {
		fmt.Printf("失败，错误信息： %v\n", err)
		return
	}

	// 切换日志文件
	postFix := "_" + now.Add(time.Hour*24).Format("2006-01-02 15-04-05")
	logPath := filepath.Join(ioconfig.GetLoggerConf().Path, logFile+postFix)
	logOut, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("创建新的日志文件 %s 失败，错误信息: %v\n", logFile, err))
		panic(err)
	}
	debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
	infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
	warnLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
	errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
	timeFlag = time.Now()

	// 判断创建时间是否超过最大允许的时间差并删除过期文件
	files, err := os.ReadDir(ioconfig.GetLoggerConf().Path)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("无法打开文件夹，错误信息: %v\n", err))
		return
	}
	for _, file := range files {
		// 获取文件信息
		filePath := filepath.Join(ioconfig.GetLoggerConf().Path, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("获取文件 %s 的信息失败，错误信息: %v\n", file.Name(), err))
			continue
		}
		// 获取文件的创建时间
		creationTime := fileInfo.ModTime()
		if now.Sub(creationTime) > maxAge {
			err := os.Remove(filePath)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("删除文件 %s 失败，错误信息: %v\n", file.Name(), err))
			}
		}
	}
}

// 给每行日志添加前缀：文件名加行号
func addPrefix() string {
	file, _, line := util.GetFileAndLineNo(3) // GetFileAndLineNo() 是第 0 层 --> addPrefix() 是第 1 层 --> Info() 是第 2 层 --> logger.Info() 的调用方是第 3 层
	arr := strings.Split(file, "/")
	if len(arr) > 5 {
		arr = arr[len(arr)-5:] // 完整的绝对路径太长，只取最末 5 级路径
	}
	return strings.Join(arr, "/") + ":" + strconv.Itoa(line) //文件名加行号
}

func Debug(format string, v ...any) {
	if logLevel <= ioconst.LoggerLevel_value["debug"] {
		checkLogRotation()
		debugLogger.Printf(addPrefix()+" "+format, v...) // format 末尾如果没有换行符会自动加上
	}
}

func Info(format string, v ...any) {
	if logLevel <= ioconst.LoggerLevel_value["info"] {
		checkLogRotation()
		infoLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Warn(format string, v ...any) {
	if logLevel <= ioconst.LoggerLevel_value["warn"] {
		checkLogRotation()
		warnLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Error(format string, v ...any) {
	if logLevel <= ioconst.LoggerLevel_value["error"] {
		checkLogRotation()
		errorLogger.Printf(addPrefix()+" "+format, v...)
	}
}

// Panicf 适用于无法恢复的错误，如文件无法打开、数据库连接失败等，在这些情况下，程序可能无法继续执行
func Panicf(format string, v ...any) {
	checkLogRotation()
	errorLogger.Panicf(addPrefix()+" "+format, v...)
}

// Fatalf 适用于严重错误，可能会导致程序执行不正确或不安全的情况
func Fatalf(format string, v ...any) {
	checkLogRotation()
	errorLogger.Panicf(addPrefix()+" "+format, v...)
}
