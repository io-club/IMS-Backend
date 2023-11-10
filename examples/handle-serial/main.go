package main

import (
	"IMS-Backend/examples/handle-serial/coordinator"
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tarm/serial"
)

var uartDevice = "/dev/ttyUSB0"
var savePath = "./data"

func main() {
	err := godotenv.Load(".env.sample")
	if err == nil {
		// 如果配置文件加载成功，就使用配置文件中的环境变量
		uartDevice = os.Getenv("UART_DEVICE")
		savePath = os.Getenv("SAVE_PATH")
	}

	// 数据位 8, 停止位 1, 校验位 None
	config := &serial.Config{
		Name:        uartDevice,
		Baud:        9600,
		ReadTimeout: 2 * time.Second,
	}
	ser, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	defer ser.Close()

	f, err := os.OpenFile(savePath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("打开日志文件失败", err)
	}
	defer f.Close()
	log.Println(config.Name)

	// before := time.Now()
	for {
		// 读取一行
		x, _, err := bufio.NewReader(ser).ReadLine()
		fmt.Printf("%v\n", string(x))
		if err != nil {
			continue
		}

		// 写入日志
		if x != nil && len(x) != 0 {
			s := time.Now().Format("2006-01-02 15:04:05")
			_, err := f.WriteString(fmt.Sprintf("[%s] %s", s, x))
			if err != nil {
				log.Println("日志写入失败", err)
			}
		}

		// 解析为结构体
		uartMsg := coordinator.DecodeUartMsgList(string(x))
		if uartMsg == nil {
			continue
		}

		// 根据传感器类型做不同的处理
		sensor := uartMsg.Sensor
		tID := uartMsg.TerminalID
		if _, ok := sensor.(*coordinator.TAndHSensor); ok {
			log.Printf("终端：%s 温度：%d 湿度：%d\n", tID, sensor.(*coordinator.TAndHSensor).Temperature, sensor.(*coordinator.TAndHSensor).Humidity)
		} else if _, ok := sensor.(*coordinator.GasSensor); ok {
			log.Printf("终端：%s 气体：%d\n", tID, sensor.(*coordinator.GasSensor).Data)
		} else if _, ok := sensor.(*coordinator.LightSensor); ok {
			log.Printf("终端：%s 光敏：%d\n", tID, sensor.(*coordinator.LightSensor).Data)
		} else {
			log.Printf("终端：%s Unknown: %s\n", tID, sensor.(*coordinator.UnknownSensor).Data)
		}

		// 计算间隔时间
		// current := time.Now()
		// interval := current.Sub(before)
		// if interval >= 0.01 {
		// 	fmt.Println(current.Sub(before))
		// } else {
		// 	fmt.Println()
		// }
	}
}
