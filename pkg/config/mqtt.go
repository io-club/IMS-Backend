package ioconfig

import (
	"log"
	"sync"
)

var (
	mqttConf *MqttConf
	mqttOnce sync.Once
)

type MqttConf struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	ClientID string `mapstructure:"client_id"`
}

func GetMqttConf() *MqttConf {
	mqttOnce.Do(func() {
		if mqttConf == nil {
			if err := V.UnmarshalKey("mqtt", &mqttConf); mqttConf == nil || err != nil {
				log.Panicf("unmarshal conf failed, err: %s \n", err)
			}
		}
	})
	return mqttConf
}
