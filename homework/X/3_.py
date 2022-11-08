from fastapi import FastAPI
from fastapi.responses import JSONResponse, HTMLResponse, FileResponse

app = FastAPI()

@app.get("/user")
def user():
    return JSONResponse(content={"msg":"get user"},
    status_code=202,
    headers={"a": "b"})

@app.get("/")
def user():
    html_content = """
    <html>
        <body><p style="color:red">Hello World</p></body>
        </html>
        """
    return HTMLResponse(content=html_content)


@app.get("/avatar")
def user():
    avatar = 'Pictures/2.jpg'
    return FileResponse(avatar)

@app.get("/2")
def user():
    avatar1 = 'Pictures/3.jpg'
    return FileResponse(avatar1)


    

