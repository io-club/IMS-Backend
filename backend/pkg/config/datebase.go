package ioconfig

import (
	"log"
	"sync"
)

var (
	dbConf *DBConf
	dbOnce sync.Once
)

type MySQLConf struct {
	GetSlaveStrategy string
	Master           string
	Slaves           []string
	Username         string
	Password         string
	Host             string
	Port             int
	DBName           string
}

type SQLiteConf struct {
	Master string
}

type DBConf struct {
	DBType string     `mapstructure:"db"`
	MySQL  MySQLConf  `mapstructure:"mysql"`
	SQLite SQLiteConf `mapstructure:"sqlite"`
}

func GetDBConf() *DBConf {
	dbOnce.Do(func() {
		if dbConf == nil {
			if err := V.UnmarshalKey("database", &dbConf); dbConf == nil || err != nil {
				log.Panicf("unmarshal conf failed, err: %s\n", err)
			}
		}
	})
	return dbConf
}
