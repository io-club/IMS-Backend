package iomqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	ioconfig "ims-server/pkg/config"
	iolog "ims-server/pkg/logger"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	iolog.Info("mqtt received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	iolog.Error("mqtt connect lost,err: %v", err)
}

type mqttClient struct {
	client mqtt.Client
}

func NewClient() (*mqttClient, error) {
	config := ioconfig.GetMqttConf()

	addr := fmt.Sprintf("tcp://%s:%s", config.Host, config.Port)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(addr)
	if config.User != "" {
		opts.SetUsername(config.User)
		if config.Password != "" {
			opts.SetPassword(config.Password)
		}
	}
	if config.ClientID != "" {
		opts.SetClientID(config.ClientID)
	}
	opts.SetDefaultPublishHandler(messagePubHandler) // Set the default message handler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		iolog.Error("mqtt connect failed, err: %v", token.Error())
		return nil, token.Error()
	}

	return &mqttClient{client: client}, nil
}
