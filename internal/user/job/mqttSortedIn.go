//接收与发布
package job

import (
    "fmt"
    mqtt "github.com/eclipse/paho.mqtt.golang"
   // "log"
    "time"
    Co "ims-server/pkg/Mqtt"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
    uartMsg := Co.DecodeUartMsgList(msg)
    //http.post
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
    fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    fmt.Printf("Connect lost: %v", err)
}



func main() {
    var broker = "broker.emqx.io"
    var port = 1883
    opts := mqtt.NewClientOptions()  //赋值为结构体
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))  //结构体调用函数
    opts.SetClientID("go_mqtt_client")
    opts.SetUsername("emqx")
    opts.SetPassword("public")
    opts.SetDefaultPublishHandler(messagePubHandler)
    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    sub(client)
    publish(client)

    client.Disconnect(250)
}

func publish(client mqtt.Client) {
    num := 10
    for i := 0; i < num; i++ {
        text := fmt.Sprintf("Message %d", i)
        token := client.Publish("topic/test", 0, false, text) //主题，qos,
        token.Wait()
        time.Sleep(time.Second)
    }
}

func sub(client mqtt.Client) {
    topic := "topic/test"
    token := client.Subscribe(topic, 1, nil)
    token.Wait()
  fmt.Printf("Subscribed to topic: %s", topic)
}