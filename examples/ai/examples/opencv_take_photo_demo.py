import cv2


def take_photo():
    # 打开摄像头
    capture = cv2.VideoCapture(0)  # 0 为电脑内置摄像头

    # 拍一张照
    ret, frame = capture.read()  # 摄像头读取，ret 为是否成功打开摄像头，true,false。frame 为视频的每一帧图像
    if not ret:
        return

    # 翻转一下
    frame = cv2.flip(frame, 1)  # 摄像头是和人对立的，将图像左右调换回来正常显示。

    # 存到文件里
    cv2.imwrite("base.png", frame)

take_photo()
