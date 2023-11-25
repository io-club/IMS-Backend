package ioconfig

import (
	iologger "ims-server/pkg/logger"
	"sync"
)

var (
	serviceConf *ServiceConf
	serviceOnce sync.Once
)

type Service struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type ServiceConf struct {
	Nms  Service `mapstructure:"nms"`
	User Service `mapstructure:"user"`
}

// The returned ServieConf will not guarantee that individual services are not empty, but it will ensure that the overall configuration is not empty. This is done to ensure that the configuration is consistent across all services.
func GetServiceConf() *ServiceConf {
	serviceOnce.Do(func() {
		if serviceConf == nil {
			if err := V.UnmarshalKey("services", &serviceConf); serviceConf == nil || err != nil {
				iologger.Panicf("unmarshal conf failed, err: %s\n", err)
			}
		}
	})
	return serviceConf
}
