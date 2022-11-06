import time
import serial
from coordinator import GasSensor, LightSensor, TAndHSensor, UnknownSensor, decode_uart_msg

def main():
    # 数据位 8, 停止位 1, 校验位 None
    with serial.Serial('/dev/ttyUSB1', 9600, timeout=2) as ser:
        print(ser.name)

        before = time.time()
        while True:

            # 读取一行
            x = ser.read_until(b"\n").decode('ascii')
            # x = ser.readline() # \0 #.decode('ascii'),

            # 解析为结构体
            uart_msg = decode_uart_msg(x)
            if uart_msg is None:
                continue

            print(f"终端：{uart_msg.terminal_id}", end=" ")

            # 根据传感器类型做不同的处理
            data = uart_msg.data
            match data:
                case TAndHSensor():
                    print(f"温度： {data.temperature} 湿度： {data.humidity}")
                    pass
                case GasSensor():
                    print("气体：", data.data)
                case LightSensor():
                    print("光敏：", data.data)
                case UnknownSensor():
                    print("Unknown:", data.data)

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
