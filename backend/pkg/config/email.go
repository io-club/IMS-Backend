package ioconfig

import (
	"log"
	"sync"
)

var (
	emailConf *EmailConf
	emailOnce sync.Once
)

type EmailConf struct {
	MailUserName string `mapstructure:"mailUserName"`
	MailPassword string `mapstructure:"mailPassword"`
	Addr         string `mapstructure:"addr"`
	Host         string `mapstructure:"host"`
}

func GetEmailConf() *EmailConf {
	emailOnce.Do(func() {
		if emailConf == nil {
			if err := V.UnmarshalKey("email", &emailConf); emailConf == nil || err != nil {
				log.Panicf("unmarshal conf failed, err: %s \n", err)
			}
		}
	})
	return emailConf
}
