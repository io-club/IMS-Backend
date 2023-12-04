package iolog

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
	"syscall"
	"time"
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger

	logLevel = 0 //The default LogLevel is 0, that is, all levels of logs are printed
)

var (
	config ioconfig.LoggerConf

	logOut    *os.File      // Log output file
	heartbeat time.Duration // Log heartbeat interval
	maxAge    time.Duration // Log retention time
	logFile   string        // Log file name

	timeFlag      time.Time    // Marks the last log switching time
	dayChangeLock sync.RWMutex // Lock when switching log files
)

func SetLogger(serviceName string) {
	var ok bool
	var err error
	value, ok := ioconfig.ServiceConfMap[serviceName]
	if !ok {
		panic("Error occurred when getting the \"service name\" configuration")
	}

	config = value.(ioconfig.Service).LoggerConf
	// 初始化 level,heartbeat,maxAge
	logLevel, ok = ioconst.LoggerLevelValue[config.Level]
	if !ok {
		panic("Error occurred when getting the \"log level\" configuration")
	}
	heartbeat = config.HeartBeat
	maxAge = config.MaxAge

	// Lock when opening the file to avoid contention
	dayChangeLock.Lock()
	var filePath string
	// Assign the file from the configuration by default
	logFile = config.FileName
	postFix := "_" + time.Now().Format("2006-01-02 15-04-05")
	filePath = filepath.Join(config.Path, logFile+postFix)
	// Try to find a file that has not exceeded the maximum age for writing
	files, err := os.ReadDir(config.Path)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Unable to open folder, error: %v\n", err))
		return
	}
	for _, file := range files {
		path := filepath.Join(config.Path, file.Name())
		fileInfo, err := os.Stat(path) // If the file is already in use, an error indicating that the file is in use will be returned
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Failed to get information for file %s, error: %v\n", file.Name(), err))
			continue
		}
		// TODO: 注意：获取文件创建时间的方法在不同的操作系统上可能会有所不同，能优化优化吗？
		// 在 window 操作系统上，使用 fileInfo.Sys().(*syscall.Win32FileAttributeData).CreationTime 来获取
		// 在 linux 操作系统上，使用 fileInfo.Sys().(*syscall.Stat_t).Ctim.Sec 来获取 (在 windows 环境下，该语句会报错，所以只能手动切换？)
		creationTime := fileInfo.Sys().(*syscall.Stat_t).Ctim.Sec
		// creationTime := fileInfo.Sys().(*syscall.Win32FileAttributeData).CreationTime.Nanoseconds() / 1e9
		creatTime := time.Unix(creationTime, 0)
		if time.Now().Sub(creatTime) <= maxAge {
			// If there is a file that has not exceeded the maximum age, override the default values
			logFile = fileInfo.Name()
			filePath = path
			break
		}
	}
	logOut, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	// Successfully opened the file, release the lock
	dayChangeLock.Unlock()

	debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
	debugLogger.SetOutput(io.MultiWriter(debugLogger.Writer(), os.Stdout)) // Synchronize output to terminal

	infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
	infoLogger.SetOutput(io.MultiWriter(infoLogger.Writer(), os.Stdout))

	warnLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
	warnLogger.SetOutput(io.MultiWriter(warnLogger.Writer(), os.Stdout))

	errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
	errorLogger.SetOutput(io.MultiWriter(errorLogger.Writer(), os.Stdout))

	timeFlag = time.Now()
	dayChangeLock = sync.RWMutex{}
}

func checkLogRotation() {
	now := time.Now()
	// Most of the time it will return directly, so it's not slow
	if now.Sub(timeFlag) <= heartbeat {
		return
	}
	// Need to switch log files, lock
	dayChangeLock.Lock()
	defer dayChangeLock.Unlock()

	err := logOut.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Failed to close log file %s, error: %v\n", logFile, err))
		return
	}

	// Switch log files
	postFix := "_" + now.Format("2006-01-02 15-04-05")
	logPath := filepath.Join(config.Path, logFile+postFix)
	logOut, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Failed to create new log file %s, error: %v\n", logFile, err))
		panic(err)
	}
	debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
	infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
	warnLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
	errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
	timeFlag = time.Now()

	// Judge whether the creation time exceeds the maximum allowable time difference and delete expired files
	files, err := os.ReadDir(config.Path)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Unable to open folder, error information: %v\n", err))
		return
	}
	for _, file := range files {
		// Get file information
		filePath := filepath.Join(config.Path, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Failed to get information for file %s, error information: %v\n", file.Name(), err))
			continue
		}
		// Get the file creation time
		creationTime := fileInfo.Sys().(*syscall.Stat_t).Ctim.Sec
		// creationTime := fileInfo.Sys().(*syscall.Win32FileAttributeData).CreationTime.Nanoseconds() / 1e9
		creatTime := time.Unix(creationTime, 0)
		if now.Sub(creatTime) > maxAge {
			err := os.Remove(filePath)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Failed to delete file %s, error information: %v\n", file.Name(), err))
			}
		}
	}
}

// Add prefix to each log line: filename and line number
func addPrefix() string {
	file, _, line := util.GetFileAndLineNo(3) // GetFileAndLineNo() is the 0th layer --> addPrefix() is the 1st layer --> Info() is the 2nd layer --> the caller of logger.Info() is the 3rd layer
	arr := strings.Split(file, "/")
	if len(arr) > 5 {
		arr = arr[len(arr)-5:] // The complete absolute path is too long, only take the last 5 path levels
	}
	return strings.Join(arr, "/") + ":" + strconv.Itoa(line) // Filename and line number
}

func Debug(format string, v ...any) {
	if logLevel <= ioconst.LoggerLevelValue["debug"] {
		checkLogRotation()
		debugLogger.Printf(addPrefix()+" "+format, v...) // If there's no newline character at the end of the format, it will be automatically added
	}
}

func Info(format string, v ...any) {
	if logLevel <= ioconst.LoggerLevelValue["info"] {
		checkLogRotation()
		infoLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Warn(format string, v ...any) {
	if logLevel <= ioconst.LoggerLevelValue["warn"] {
		checkLogRotation()
		warnLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Error(format string, v ...any) {
	if logLevel <= ioconst.LoggerLevelValue["error"] {
		checkLogRotation()
		errorLogger.Printf(addPrefix()+" "+format, v...)
	}
}

// Panicf is applicable to irrecoverable errors, such as file unable to open, database connection failure, etc. In these cases, the program may not be able to continue to run
func Panicf(format string, v ...interface{}) {
	checkLogRotation()
	errorLogger.Panicf(addPrefix()+" "+format, v...)
}

// Fatalf is applicable to serious errors that may cause the program to run incorrectly or unsafely
func Fatalf(format string, v ...interface{}) {
	checkLogRotation()
	errorLogger.Fatalf(addPrefix()+" "+format, v...)
}
