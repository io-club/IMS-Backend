package iomqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	iolog "ims-server/pkg/logger"
)

func Sub(client mqtt.Client, topic string, handler mqtt.MessageHandler) error {
	if token := client.Subscribe(topic, 1, handler); token.Wait() && token.Error() != nil {
		iolog.Error("mqtt subscribe failed, err: %v", token.Error())
		return token.Error()
	}
	iolog.Info("mqtt subscribe success, topic: %s", topic)
	return nil
}

func Unsub(client mqtt.Client, topic string) {
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		iolog.Error("mqtt unsubscribe failed, err: %v", token.Error())
	}
}

func SetRetained(client mqtt.Client, topic string, text string) {
	token := client.Publish(topic, 0, true, text)
	token.Wait()
}

func Publish(client mqtt.Client, topic string, text string, qos ...byte) {
	quality := byte(0)
	if qos != nil {
		quality = qos[0]
	}
	token := client.Publish(topic, quality, false, text)
	token.Wait()
}

func MPublish(client mqtt.Client, topic string, texts []string, qos ...byte) {
	quality := byte(0)
	if qos != nil {
		quality = qos[0]
	}
	for _, t := range texts {
		token := client.Publish(topic, quality, false, t)
		token.Wait()
	}
}
