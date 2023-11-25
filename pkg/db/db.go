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
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
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
