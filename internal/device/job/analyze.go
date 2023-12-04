package job

import (
	"fmt"
	iolog "ims-server/pkg/logger"
	"ims-server/pkg/mqtt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)


var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	iolog.Info("mqtt received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	fmt.Println((string(msg.Payload())))
}

func ClientSub()(){
	client, err := iomqtt.NewClient()
	if err != nil {
		return
	}

	if err := client.Sub("abc", messagePubHandler); err != nil {
		return
	}
}

func ClientPub (msg string)() {
	client, err := iomqtt.NewClient()
	if err != nil {
		return
	}
	client.Publish("abc", msg, 0)
}