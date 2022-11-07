import time
import serial
from coordinator import GasSensor, LightSensor, TAndHSensor, UnknownSensor, decode_uart_msg
import json


def main():
    # 数据位 8, 停止位 1, 校验位 None
    with serial.Serial('/dev/ttyUSB1', 9600, timeout=2) as ser, open("./data", "a+") as f:
        print(ser.name)

        before = time.time()
        while True:

            # 读取一行
            x = ser.read_until(b"\n").decode('ascii')
            # x = ser.readline() # \0 #.decode('ascii'),

            # 写入日志
            if x != None and len(x) != 0:
                s = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())
                f.write(f"[{s}] {x}")

            # 解析为结构体
            uart_msg = decode_uart_msg(x)
            if uart_msg is None:
                continue

            print(f"终端：{uart_msg.terminal_id}", end=" ")

            # 根据传感器类型做不同的处理
            sensor = uart_msg.sensor
            if isinstance(sensor, TAndHSensor):
                print(f"温度： {sensor.temperature} 湿度： {sensor.humidity}")
                pass
            elif isinstance(sensor, GasSensor):
                print("气体：", sensor.data)
            elif isinstance(sensor, LightSensor):
                print("光敏：", sensor.data)
            else:
                print("Unknown:", sensor.data)

            # # 计算间隔时间
            # current = time.time()
            # interval = current - before
            # if interval >= 0.01:
            #     print(current - before)
            # else:
            #     print()
            # # 重新计时
            # before = current


def test():
    # match int("1"):
    #     case 1:
    #         print("num")
    #     case "1":
    #         print("str")
    pass


if __name__ == "__main__":
    main()
    # test()
