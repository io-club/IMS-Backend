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
	if config.User != "" && config.Password != "" {
		opts.SetUsername(config.User)
		opts.SetPassword(config.Password)
	}
	if config.ClientID != "" {
		opts.SetClientID(config.ClientID)
	}
	opts.SetDefaultPublishHandler(messagePubHandler) // 设置默认的接收函数
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		iolog.Error("mqtt connect failed, err: %v", token.Error())
		return nil, token.Error()
	}

	return &mqttClient{client: client}, nil
}
