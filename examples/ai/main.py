from fdlite.render import Annotation, Point, Line, RectOrOval, FilledRectOrOval
from PIL.Image import Image as PILImage
from PIL import ImageDraw
from typing import Sequence, Union
from time import sleep
from typing import Union

import cv2
from fastapi import FastAPI
from fastapi.responses import StreamingResponse
import numpy as np
from PIL import Image


app = FastAPI()


@app.get("/")
def read_root():
    return {"Hello": "World"}


@app.get("/items/{item_id}")
def read_item(item_id: int, q: Union[str, None] = None):
    return {"item_id": item_id, "q": q}


@app.get("/video")
def video():
    def iterfile():  # (1)
        capture = cv2.VideoCapture(0)  # 0 为电脑内置摄像头
        while(True):
            sleep(1.0 / 30)

            ret, frame = capture.read()  # 摄像头读取，ret 为是否成功打开摄像头，true,false。frame 为视频的每一帧图像
            if not ret:
                break

            # normal CV2
            frame = cv2.flip(frame, 1)  # 摄像头是和人对立的，将图像左右调换回来正常显示。

            from fdlite import FaceDetection, FaceDetectionModel
            from fdlite.render import Colors, detections_to_render_data

            detect_faces = FaceDetection(
                model_type=FaceDetectionModel.BACK_CAMERA)
            faces = detect_faces(frame)
            render_data = detections_to_render_data(
                faces, bounds_color=Colors.GREEN)

            # renderd PIL
            from PIL.Image import Image as PILImage
            frame: PILImage = render_to_image(render_data, toImgPIL(frame))

            # renderd CV2
            frame = toImgOpenCV(frame)
            # encoded cv2
            ret, buffer = cv2.imencode('.jpg', frame)
            # bytes cv2
            frame = buffer.tobytes()

            yield (b'--frame\r\n'
                   b'Content-Type: image/jpeg\r\n\r\n' + frame + b'\r\n')  # concat frame one by one and show result

    return StreamingResponse(iterfile(), media_type="multipart/x-mixed-replace; boundary=frame")


def toImgOpenCV(imgPIL):  # Conver imgPIL to imgOpenCV
    i = np.array(imgPIL)  # After mapping from PIL to numpy : [R,G,B,A]
    # numpy Image Channel system: [B,G,R,A]
    red = i[:, :, 0].copy()
    i[:, :, 0] = i[:, :, 2].copy()
    i[:, :, 2] = red
    return i


def toImgPIL(imgOpenCV):
    return Image.fromarray(cv2.cvtColor(imgOpenCV, cv2.COLOR_BGR2RGB))


croped_base = cv2.imread("./croped_base.jpg")

# 这里
def render_to_image(
    annotations: Sequence[Annotation],
    image: PILImage,
    blend: bool = False
) -> PILImage:
    draw = ImageDraw.Draw(image, mode='RGBA' if blend else 'RGB')

    for annotation in annotations:
        if annotation.normalized_positions:
            scaled = annotation.scaled(image.size)
        else:
            scaled = annotation
        if not len(scaled.data):
            continue
        thickness = int(scaled.thickness)
        color = scaled.color
        for item in scaled.data:
            if isinstance(item, Point):
                w = max(thickness // 2, 1)
                rc = [item.x - w, item.y - w, item.x + w, item.y + w]
                draw.rectangle(rc, fill=color.as_tuple, outline=color.as_tuple)
            elif isinstance(item, Line):
                coords = [item.x_start, item.y_start, item.x_end, item.y_end]
                draw.line(coords, fill=color.as_tuple, width=thickness)
            elif isinstance(item, RectOrOval):
                rc = [item.left, item.top, item.right, item.bottom]

                croped = image.crop(rc).save("./croped.jpg")
                # res = get_similar("./croped.jpg", "./croped_base.jpg")
                # print(res)

                from deepface import DeepFace
                df = DeepFace.verify(img1_path="./croped.jpg",
                                     img2_path="./croped_base.jpg", enforce_detection=False)
                print(df)
                # print(df["distance"])
                # print(1.0 / (1.0 + df["distance"]))
                # print()

                if item.oval:
                    draw.ellipse(rc, outline=color.as_tuple, width=thickness)
                else:
                    draw.rectangle(rc, outline=color.as_tuple, width=thickness)

            elif isinstance(item, FilledRectOrOval):
                rgb = color.as_tuple
                rect, fill = item.rect, item.fill.as_tuple
                rc = [rect.left, rect.top, rect.right, rect.bottom]
                if rect.oval:
                    draw.ellipse(rc, fill=fill, outline=rgb, width=thickness)
                else:
                    draw.rectangle(rc, fill=fill, outline=rgb, width=thickness)

            else:
                # don't know how to render this
                pass
    return image
