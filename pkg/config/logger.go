package ioconfig

import (
	iologger "ims-server/pkg/logger"
	"strings"
	"sync"
	"time"
)

var (
	loggerConf *LoggerConf
	loggerOnce sync.Once
)

type LoggerConf struct {
	Level     string        `mapstructure:"level"`
	Path      string        `mapstructure:"path"`
	FileName  string        `mapstructure:"file_name"`
	HeartBeat time.Duration `mapstructure:"heartbeat"`
	MaxAge    time.Duration `mapstructure:"max_age"`
}

func GetLoggerConf() *LoggerConf {
	loggerOnce.Do(func() {
		if loggerConf == nil {
			if err := V.UnmarshalKey("logger", &loggerConf); loggerConf == nil || err != nil {
				iologger.Panicf("unmarshal conf failed, err: %s\n", err)
			}
			loggerConf.Level = strings.ToLower(loggerConf.Level)
			loggerConf.Path = RootPath + loggerConf.Path
			loggerConf.HeartBeat = loggerConf.HeartBeat * time.Minute
			loggerConf.MaxAge = loggerConf.MaxAge * time.Hour
		}
	})
	return loggerConf
}
