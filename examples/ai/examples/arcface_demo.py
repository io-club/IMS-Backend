from arcface import ArcFace
from numpy import number

face_rec = ArcFace.ArcFace()


# 比较人脸相似度
def get_similar(path1: str, path2: str) -> number:
    emb1 = face_rec.calc_emb(path1)
    emb2 = face_rec.calc_emb(path2)
    res = face_rec.get_distance_embeddings(emb1, emb2)
    return 1.0 / (res + 1)


if __name__ == "__main__":
    a = get_similar("../../../../../../Downloads/fb001f5ee4b1d6becf367ef10b2574fe.jpeg",
                    "../../../../../../Downloads/5c38aa31771f712ec824b3850c42b1a9.jpeg")
    print(a)

    a = get_similar("./test1.png", "./test2.png")
    print(a)
