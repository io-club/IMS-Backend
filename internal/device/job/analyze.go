package job

import (
	"fmt"
	iolog "ims-server/pkg/logger"
	"ims-server/pkg/mqtt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)


var manualMessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	iolog.Info("mqtt received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	fmt.Println(DecodeUartMsgList(string(msg.Payload())))
}


func Client_test (topic ,msg string){
	client, err := iomqtt.NewClient()
	if err != nil {
		iolog.Warn("err:%v",err)
		return
	}

	if err := client.Sub(topic, manualMessagePubHandler); err != nil {
		iolog.Warn("err:%v",err)
		return
	}

	
	for i := 0; i < 3; i++ {
		client.Publish(topic, msg)
	}

	time.Sleep(10*time.Second)

	client.Client.Disconnect(250)
}