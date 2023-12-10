package job

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

func (s SensorType) ToString() string {
	switch s {
	case SensorTypeTAndH:
		return "TAndH"
	case SensorTypeGas:
		return "Gas"
	case SensorTypeLight:
		return "Light"
	default:
		return "Unknown"
	}
}

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

type UartMsg struct {
	TerminalID string
	Sensor     interface{} //以上之一的传感器类型
}

func DecodeUartMsgList(msg string) (*SensorType, *UartMsg) {
	if len(msg) == 0 {
		return nil, nil
	}
	res := UartMsg{}
	msgParts := strings.Split(msg, ":")
	//TerminalID:sensor:Date1 Date2 ...
	res.TerminalID = msgParts[0]
	sensor, err := strconv.Atoi(msgParts[1])
	if err != nil {
		return nil, nil
	}
	Data := msgParts[2]

	sensorType := SensorType(sensor)
	switch sensorType {
	case SensorTypeTAndH:
		parts := strings.Split(Data, " ")
		temperature, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, nil
		}
		humidity, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, nil
		}
		res.Sensor = &TAndHSensor{
			Temperature: temperature,
			Humidity:    humidity,
		}
	case SensorTypeGas:
		dataInt, err := strconv.Atoi(Data)
		if err != nil {
			return nil, nil
		}
		res.Sensor = &GasSensor{
			Data: dataInt,
		}
	case SensorTypeLight:
		dataInt, err := strconv.Atoi(Data)
		if err != nil {
			return nil, nil
		}
		res.Sensor = &LightSensor{
			Data: dataInt,
		}
	default:
		res.Sensor = &UnknownSensor{
			Data: Data,
		}
	}
	return &sensorType, &res
}
