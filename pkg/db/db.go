package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	ioconfig "ims-server/pkg/config"
	ioconst "ims-server/pkg/consts"
	"log"
	"sync"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func NewDB() *gorm.DB {
	dbOnce.Do(func() {
		conf := ioconfig.GetDBConf()
		switch conf.DBType {
		case ioconst.DBTypeMySQL.String():
			db = NewMySQLProxy()
		case ioconst.DBTypeSQLite.String():
			db = NewSQLiteProxy()
		default:
			log.Panicf("db type %s not support", conf.DBType)
		}
	})
	return db
}

func NewMySQLProxy() *gorm.DB {
	conf := ioconfig.GetDBConf().MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username, conf.Password, conf.Host, conf.Port, conf.DBName,
	)
	proxy, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       dsn,
			DefaultStringSize:         256,   // Default length of string type fields
			DontSupportRenameIndex:    true,  // Use the method of deleting and recreating when renaming the index, the database before MySQL 5.7 and MariaDB does not support renaming indexes
			SkipInitializeWithVersion: false, // Automatically configure according to the current MySQL version
		}),
		&gorm.Config{
			PrepareStmt: true,
		},
	)
	if err != nil {
		log.Panicf("db connect fail,err:%v", err)
	}
	return proxy
}

func NewSQLiteProxy() *gorm.DB {
	conf := ioconfig.GetDBConf().SQLite
	proxy, err := gorm.Open(sqlite.Open(conf.Master), &gorm.Config{PrepareStmt: true})
	if err != nil {
		log.Panicf("db connect fail,err:%v", err)
	}
	return proxy
}
