import uvicorn
from fastapi import FastAPI

app = FastAPI()
#注释

#添加首页
@app.get("/")
def index():
    """这是首页"""
    return "This is Home Page."

@app.get("/users")
def users():
    """获取所有用户信息"""
    return {"msg": "Get all Users", "code": 2001}

@app.get("/projects")
def projects():
    return ["项目1", "项目2", "项目3"]
if __name__ == '__main__':
    uvicorn.run(app)

