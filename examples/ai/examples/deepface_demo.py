from deepface import DeepFace

# df = DeepFace.verify(img1_path="../../../../../Downloads/5c38aa31771f712ec824b3850c42b1a9.jpeg",
#                      img2_path="../../../../../Downloads/fb001f5ee4b1d6becf367ef10b2574fe.jpeg")
# print(df)
# print(df["distance"])
# print(1.0 / (1.0 + df["distance"]))

# df = DeepFace.verify(img1_path="../../../../../Downloads/5c38aa31771f712ec824b3850c42b1a9.jpeg",
#                      img2_path="../../../../../Downloads/5c38aa31771f712ec824b3850c42b1a9.jpeg")
# print(df)
# print(df["distance"])
# print(1.0 / (1.0 + df["distance"]))

# obj = DeepFace.analyze(img_path="../../../../../Downloads/5c38aa31771f712ec824b3850c42b1a9.jpeg",
#                        actions=['age', 'gender', 'race', 'emotion']
#                        )

# print(obj)

# 使用 deepface 比较人脸相似度
df = DeepFace.verify(img1_path="./croped.jpg",
                     img2_path="./croped_base.jpg", enforce_detection=False)
print(df)

obj = DeepFace.analyze(img_path="./test.png",
                       actions=[
                           'age',
                           'gender',
                           # 'race',
                           'emotion'],
                       enforce_detection=False
                       )

print(obj)
