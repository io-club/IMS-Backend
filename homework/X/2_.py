from fastapi import FastAPI, Header, Body, Form
app = FastAPI()

# @app.post("/login")
# def login():
#     return {"msg": "login success"}

# @app.api_route("/login",methods=("GET", "POST", "PUT"))
# def login():
#     return {"msg": "login success"}


@app.get("/user")
def user(id,token=Header(None)):
    return {"id":id, "token": token}

@app.post("/login")
def login(username=Form(None), password=Form(None)):
    return {"data":{"username": username, "password": password}}


