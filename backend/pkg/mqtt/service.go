package iomqtt

import (
	iolog "ims-server/pkg/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (mc *mqttClient) Sub(topic string, handler mqtt.MessageHandler) error {
	if token := mc.Client.Subscribe(topic, 1, handler); token.Wait() && token.Error() != nil {
		iolog.Error("mqtt subscribe failed, err: %v", token.Error())
		return token.Error()
	}
	iolog.Info("mqtt subscribe success, topic: %s", topic)
	return nil
}

func (mc *mqttClient) Unsub(topic string) {
	if token := mc.Client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		iolog.Error("mqtt unsubscribe failed, err: %v", token.Error())
	}
}

func (mc *mqttClient) SetRetained(topic string, text string) {
	token := mc.Client.Publish(topic, 0, true, text)
	token.Wait()
}

func (mc *mqttClient) Publish(topic string, text string, qos ...byte) {
	quality := byte(0)
	if qos != nil {
		quality = qos[0]
	}
	token := mc.Client.Publish(topic, quality, false, text)
	token.Wait()
}

func (mc *mqttClient) MPublish(topic string, texts []string, qos ...byte) {
	quality := byte(0)
	if qos != nil {
		quality = qos[0]
	}
	for _, t := range texts {
		token := mc.Client.Publish(topic, quality, false, t)
		token.Wait()
	}
}
