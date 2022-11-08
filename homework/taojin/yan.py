import uvicorn
from fastapi import FastAPI, Header, Body, Form, Request
from fastapi.responses import JSONResponse, HTMLResponse, FileResponse
from fastapi.templating import Jinja2Templates

app = FastAPI()
template = Jinja2Templates("pages")


@app.get("/")
def index():

    return "This is home page"


@app.get("/users")
def users():
    return {"msg": "get all users", "code": 2001}


@app.get("/projects")
def projects():
    return ["项目1", "项目二", "项目三"]


if __name__ == '__main__':
    uvicorn.run(app)


@app.post("/login")
def login():
    return {"msg": "login success"}


@app.api_route("/login", methods=("GET", "POST", "PUT"))
def login():
    return {"msg": "login success"}


@app.post("/login")
def login(username=Form(None), password=Form(None)):
    return {"data": {"username": username, "password": password}}


@app.get("/user")
def user():
    return JSONResponse(content={"msg": "get user"},
                        status_code=202,
                        headers={"a": "b"}
                        )


@app.get("/")
def user(username, req: Request):
    todos = ["写日记", "看电影"]

    return template.TemplateResponse("index.html", context={"request": req, "todos": todos})


@app.get("/avatar")
def user():
    avatar = 'static\QQ图片20221024075459.jpg'
    return FileResponse(avatar)


@app.get("/")
def user(username, req: Request):

    return template.TemplateResponse("index.html", context={"request": req, "name": username})
