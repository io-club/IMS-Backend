from typing import Union
from enum import Enum
from pydantic import BaseModel


class SensorType(Enum):
    """将 int 转为传感器类型"""

    TAndH = 1
    Gas = 2
    Light = 3


class TAndHSensor(BaseModel):
    """温湿度传感器"""

    temperature: int
    humidity: int


class GasSensor(BaseModel):
    """气体传感器"""

    data: int


class LightSensor(BaseModel):
    """光照强度传感器"""

    data: int


class UnknownSensor(BaseModel):
    """未知传感器"""

    data: str


class UartMsg:
    """串口返回数据对应的类"""

    terminal_id: str
    sensor: Union[TAndHSensor, GasSensor, LightSensor, UnknownSensor]


def decode_uart_msg(msg: str) -> Union[UartMsg, None]:
    """将串口数据解析为类"""

    # 没数据时，continue
    if len(msg) == 0:
        return None

    res = UartMsg()

    try:
        msg_parts = msg.split(":")

        res.terminal_id = msg_parts[0]

        data = msg_parts[2]
        sensor_type = SensorType(int(msg_parts[1]))
        if sensor_type == SensorType.TAndH:
            [temperature, humidity] = data.split()
            res.sensor = TAndHSensor(
                temperature=int(temperature), humidity=int(humidity)
            )
        elif sensor_type == SensorType.Gas:
            res.sensor = GasSensor(data=int(data))
        elif sensor_type == SensorType.Light:
            res.sensor = LightSensor(data=int(data))
        else:
            res.sensor = UnknownSensor(data=data)
        return res
    except:
        return None
