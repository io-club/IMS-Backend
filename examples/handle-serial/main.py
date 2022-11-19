import logging
import time

import serial
from coordinator import GasSensor, LightSensor, TAndHSensor, decode_uart_msg
from dotenv import dotenv_values

# 加载环境变量
config = dotenv_values(".env")

# 设置日志基本
LOGLEVEL = config.get("LOG_LEVEL", "INFO").upper()
logging.basicConfig(level=LOGLEVEL)


def main():
    logging.info(f"日志等级：{LOGLEVEL}")
    # 数据位 8, 停止位 1, 校验位 None
    with serial.Serial(
        config.get("UART_DEVICE", "/dev/ttyUSB0"), 9600, timeout=2
    ) as ser, open(config.get("SAVE_PATH", "./data"), "a+") as f:
        logging.info(ser.name)

        before = time.time()
        while True:

            # 读取一行
            try:
                x = ser.read_until(b"\n").decode("ascii")
                # x = ser.readline() # \0 #.decode('ascii'),
            except:
                continue

            # 写入日志
            if x != None and len(x) != 0:
                s = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())
                f.write(f"[{s}] {x}")

            # 解析为结构体
            uart_msg = decode_uart_msg(x)
            if uart_msg is None:
                continue

            # 根据传感器类型做不同的处理
            sensor = uart_msg.sensor
            t_id = uart_msg.terminal_id
            if isinstance(sensor, TAndHSensor):
                logging.debug(
                    f"终端：{t_id} 温度： {sensor.temperature} 湿度： {sensor.humidity}"
                )
                pass
            elif isinstance(sensor, GasSensor):
                logging.debug(f"终端：{t_id} 气体：{sensor.data}")
            elif isinstance(sensor, LightSensor):
                logging.debug(f"终端：{t_id} 光敏：{sensor.data}")
            else:
                logging.debug(f"终端：{t_id} Unknown: {sensor.data}")

            # # 计算间隔时间
            # current = time.time()
            # interval = current - before
            # if interval >= 0.01:
            #     print(current - before)
            # else:
            #     print()
            # # 重新计时
            # before = current


if __name__ == "__main__":
    main()
