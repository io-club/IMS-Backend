package coordinator

import (
	"strconv"
	"strings"
)

type SensorType int // 将 int 转为传感器类型

const (
	SensorTypeTAndH SensorType = 1
	SensorTypeGas   SensorType = 2
	SensorTypeLight SensorType = 3
)

type TAndHSensor struct { // 温湿度传感器
	Temperature int
	Humidity    int
}

type GasSensor struct { // 气体传感器
	Data int
}

type LightSensor struct { // 光照传感器
	Data int
}

type UnknownSensor struct { // 未知传感器
	Data string
}

type UartMsg struct { // 串口返回数据对应的类
	TerminalID string
	// 类型应为 TAndHSensor，GasSensor，LightSensor，UnknownSensor 之一
	Sensor     interface{}
}

// 将串口数据解析为类
func DecodeUartMsgList(msg string) *UartMsg {
	if len(msg) == 0 {
		return nil
	}
	res := UartMsg{}
	msgParts := strings.Split(msg, ":")
	res.TerminalID = msgParts[0]
	data := msgParts[2]
	sensor, err := strconv.Atoi(msgParts[1])
	if err != nil {
		return nil
	}
	sensorType := SensorType(sensor)
	switch sensorType {
	case SensorTypeTAndH:
		parts := strings.Split(data, "  ")
		temperature, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil
		}
		humidity, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil
		}
		res.Sensor = &TAndHSensor{
			Temperature: temperature,
			Humidity:    humidity,
		}
	case SensorTypeGas:
		dataInt, err := strconv.Atoi(data)
		if err != nil {
			return nil
		}
		res.Sensor = &GasSensor{
			Data: dataInt,
		}
	case SensorTypeLight:
		dataInt, err := strconv.Atoi(data)
		if err != nil {
			return nil
		}
		res.Sensor = &LightSensor{
			Data: dataInt,
		}
	default:
		res.Sensor = &UnknownSensor{
			Data: data,
		}
	}
	return &res
}
