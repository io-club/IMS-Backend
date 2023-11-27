package ioconfig

import (
	"ims-server/pkg/util"
	"log"
	"sync"
	"time"
)

var (
	serviceConf *ServiceConf
	serviceOnce sync.Once

	ServiceConfMap map[string]interface{}
)

type LoggerConf struct {
	Level     string        `mapstructure:"level"`
	Path      string        `mapstructure:"path"`
	FileName  string        `mapstructure:"file_name"`
	HeartBeat time.Duration `mapstructure:"heartbeat"`
	MaxAge    time.Duration `mapstructure:"max_age"`
}

type Service struct {
	Name       string `mapstructure:"name"`
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	LoggerConf `mapstructure:"logger"`
}

type ServiceConf struct {
	Nms  Service `mapstructure:"nms"`
	User Service `mapstructure:"user"`
}

// GetServiceConf The returned ServieConf will not guarantee that individual services are not empty, but it will ensure that the overall configuration is not empty. This is done to ensure that the configuration is consistent across all services.
func GetServiceConf() *ServiceConf {
	serviceOnce.Do(func() {
		if serviceConf == nil {
			if err := V.UnmarshalKey("services", &serviceConf); serviceConf == nil || err != nil {
				log.Panicf("unmarshal conf failed, err: %s\n", err)
			}
			ServiceConfMap = util.StructToMap(serviceConf)
		}
	})
	return serviceConf
}
