package main

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"ims-server/microservices/device/internal/dal/repo"
	"ims-server/microservices/device/internal/job"
	ioconfig "ims-server/pkg/config"
	iolog "ims-server/pkg/logger"
	"ims-server/pkg/mqtt"
	"os"
	"os/signal"
)

const DefaultTopic = "io-club-ims"

var manualMessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	typeof, rawData := job.DecodeUartMsgList(string(msg.Payload()))
	if typeof == nil || rawData == nil {
		iolog.Info("mqtt received message failed\n")
	}
	err := repo.NewDataRepo().CreateFromRawData(context.Background(), *typeof, rawData)
	if err != nil {
		iolog.Warn("create data failed, err: %v", err)
	}
	iolog.Info("mqtt received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func main() {
	config := ioconfig.GetServiceConf().Device
	serviceName := config.Name
	iolog.SetLogger(serviceName)

	client, err := iomqtt.NewClient()
	if err != nil {
		iolog.Panicf("NewClient Error:%v", err)
	}

	if err := client.Sub(DefaultTopic, manualMessagePubHandler); err != nil {
		iolog.Warn("Subscription Error:%v", err)
		return
	}
	// TODO：测试用，之后删除
	//for {
	//	client.Publish(DefaultTopic, "1:1:1 1")
	//	fmt.Printf("publish: 1:1:1 1\n")
	//	time.Sleep(30 * time.Second)
	//}
	// Wait for the interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
