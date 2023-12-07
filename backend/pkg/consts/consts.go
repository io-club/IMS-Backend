package ioconsts

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

// DBType 数据库类型
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

// LoggerLevel 日志级别
type LoggerLevel int

const (
	LoggerLevelDebug LoggerLevel = 0
	LoggerLevelInfo  LoggerLevel = 1
	LoggerLevelWarn  LoggerLevel = 2
	LoggerLevelError LoggerLevel = 3
)

var LoggerLevelName = map[int]string{
	0: "debug",
	1: "info",
	2: "warn",
	3: "error",
}

var LoggerLevelValue = map[string]int{
	"debug": 0,
	"info":  1,
	"warn":  2,
	"error": 3,
}

// UserType 用户类型
type UserType string

const (
	UserTypeOutsiders UserType = "outsiders" // 外部人员
	UserTypeInsider   UserType = "insider"   // 内部人员
	UserTypeAdmin     UserType = "admin"     // 管理员
)

func (u UserType) Valid() bool {
	switch u {
	case UserTypeOutsiders, UserTypeInsider, UserTypeAdmin:
		return true
	default:
		return false
	}
}

func (u UserType) String() string {
	return string(u)
}

// AccountStatus 账号状态
type AccountStatus string

const (
	AccountStatusNormal   AccountStatus = "normal"  // 正常
	AccountStatusDisabled AccountStatus = "disable" // 禁用
)

func (a AccountStatus) Valid() bool {
	switch a {
	case AccountStatusNormal, AccountStatusDisabled:
	default:
		return false
	}
	return true
}

func (a AccountStatus) String() string {
	return string(a)
}
