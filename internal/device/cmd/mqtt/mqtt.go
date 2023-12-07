package main

import (
	// "fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"ims-server/internal/device/job"
	iolog "ims-server/pkg/logger"
	"ims-server/pkg/mqtt"
)

var manualMessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	iolog.Info("mqtt received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	ParsedDate := job.DecodeUartMsgList(string(msg.Payload()))
}

var topic string = "ioMqtt"

func main() {
	clientSub, err := iomqtt.NewClient()
	if err != nil {
		iolog.Warn("err:%v", err)
		return
	}

	for {
		if err := clientSub.Sub(topic, manualMessagePubHandler); err != nil {
			iolog.Warn("err:%v", err)
			return
		}
	}
}
