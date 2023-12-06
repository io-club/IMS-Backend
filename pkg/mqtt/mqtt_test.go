package iomqtt

import (
	ioconfig "ims-server/pkg/config"
	iologger "ims-server/pkg/logger"
	"log"
	"testing"

)

func TestMqtt(t *testing.T) {
	config := ioconfig.GetServiceConf().Device
	serviceName := config.Name

	iologger.SetLogger(serviceName)

	client, err := NewClient()
	if err != nil {
		log.Printf("mqtt connect failed, err: %v", err)
		return
	}

	go func() {
		err := client.Sub("ims",nil)
		if err != nil {
			log.Printf("mqtt subscribe failed, err: %v", err)
		}
	}()

	for {
		client.Publish("ims","test")
	}
}
