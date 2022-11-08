from time import sleep
import cv2
from arcface_demo import get_similar

base = cv2.imread("./base.png")


def video_demo():
    capture = cv2.VideoCapture(0)  # 0 为电脑内置摄像头
    while(True):
        ret, frame = capture.read()  # 摄像头读取，ret 为是否成功打开摄像头，true,false。frame 为视频的每一帧图像
        if not ret:
            break

        frame = cv2.flip(frame, 1)  # 摄像头是和人对立的，将图像左右调换回来正常显示。
        res = get_similar(base, frame)
        print(res)
        cv2.imwrite("test.png", frame)
        sleep(1.0 / 5)


video_demo()
cv2.destroyAllWindows()
