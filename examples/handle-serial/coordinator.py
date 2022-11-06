
from typing import Union
from enum import Enum
from pydantic import BaseModel

# 将 int 转为传感器类型
class SensorType(Enum):
    TAndH = 1
    Gas = 2
    Light = 3

# 每个传感器对应的类
class TAndHSensor(BaseModel):
    temperature: int
    humidity: int

class GasSensor(BaseModel):
    data: int

class LightSensor(BaseModel):
    data: int

class UnknownSensor(BaseModel):
    data: str

# 串口数据对应的类
class UartMsg:
    terminal_id: str
    data: Union[TAndHSensor, GasSensor, LightSensor, UnknownSensor]

# 将串口数据解析为类
def decode_uart_msg(msg: str) -> Union[UartMsg, None]:
    # 没数据时，continue
    if len(msg) == 0:
        return None

    res = UartMsg()

    msg_parts = msg.split(':')

    res.terminal_id = msg_parts[0]

    data = msg_parts[2]
    match SensorType(int(msg_parts[1])):
        case SensorType.TAndH:
            [temperature, humidity] = data.split()
            res.data = TAndHSensor(temperature= int(temperature), humidity= int(humidity))
        case SensorType.Gas:
            res.data = GasSensor(data=int(data))
        case SensorType.Light:
            res.data = LightSensor(data=int(data))
        case _:
            res.data = UnknownSensor(data=data)
    return res
