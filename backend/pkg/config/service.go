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
	Nms    Service `mapstructure:"nms"`
	User   Service `mapstructure:"user"`
	Device Service `mapstructure:"device"`
	Work   Service `mapstructure:"work"`
}

// GetServiceConf The returned ServieConf will not guarantee that individual services are not empty, but it will ensure that the overall configuration is not empty. This is done to ensure that the configuration is consistent across all services.
func GetServiceConf() *ServiceConf {
	serviceOnce.Do(func() {
		if serviceConf == nil {
			if err := V.UnmarshalKey("services", &serviceConf); serviceConf == nil || err != nil {
				log.Panicf("unmarshal conf failed, err: %s\n", err)
			}
			// Convert time units
			serviceConf.Nms.LoggerConf.HeartBeat = serviceConf.Nms.LoggerConf.HeartBeat * time.Hour
			serviceConf.Nms.LoggerConf.MaxAge = serviceConf.Nms.LoggerConf.MaxAge * time.Hour

			serviceConf.User.LoggerConf.HeartBeat = serviceConf.User.LoggerConf.HeartBeat * time.Hour
			serviceConf.User.LoggerConf.MaxAge = serviceConf.User.LoggerConf.MaxAge * time.Hour

			serviceConf.Device.LoggerConf.HeartBeat = serviceConf.Device.LoggerConf.HeartBeat * time.Hour
			serviceConf.Device.LoggerConf.MaxAge = serviceConf.Device.LoggerConf.MaxAge * time.Hour

			serviceConf.Work.LoggerConf.HeartBeat = serviceConf.Work.LoggerConf.HeartBeat * time.Hour
			serviceConf.Work.LoggerConf.MaxAge = serviceConf.Work.LoggerConf.MaxAge * time.Hour
			// Convert the struct to a map
			ServiceConfMap = util.StructToMap(serviceConf)
		}
	})
	return serviceConf
}
