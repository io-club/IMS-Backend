from fdlite import FaceDetection, FaceDetectionModel
from fdlite.render import Colors, detections_to_render_data, render_to_image
from PIL import Image

# 人脸检测

image = Image.open('example1.jpg')

detect_faces = FaceDetection(model_type=FaceDetectionModel.BACK_CAMERA)
faces = detect_faces(image)

if not len(faces):
    print('no faces detected :(')
else:
    render_data = detections_to_render_data(faces, bounds_color=Colors.GREEN)
    render_to_image(render_data, image).show()
