package ioconst

type Mode string

const (
	// ModeRelease for production
	ModeRelease Mode = "release"
	// ModeTest for unit testing
	ModeTest Mode = "test"
	// ModeDebug for development and joint debugging tests
	ModeDebug Mode = "debug"
)

func (m Mode) String() string {
	return string(m)
}

// 数据库类型
type DBType string

const (
	DBTypeMySQL  DBType = "mysql"
	DBTypeSQLite DBType = "sqlite"
)

func (d DBType) String() string {
	return string(d)
}

func (d DBType) Valid() bool {
	switch d {
	case DBTypeMySQL, DBTypeSQLite:
		return false
	default:
		return true
	}
}

// 日志级别
type LoggerLevel int

const (
	LoggerLevelDebug LoggerLevel = 0
	LoggerLevelInfo  LoggerLevel = 1
	LoggerLevelWarn  LoggerLevel = 2
	LoggerLevelError LoggerLevel = 3
)

var LoggerLevel_name = map[int]string{
	0: "debug",
	1: "info",
	2: "warn",
	3: "error",
}

var LoggerLevel_value = map[string]int{
	"debug": 0,
	"info":  1,
	"warn":  2,
	"error": 3,
}
