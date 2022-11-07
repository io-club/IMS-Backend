from fastapi import FastAPI
from pymongo import MongoClient
from routes import router as book_router

app = FastAPI()

@app.on_event("startup")
def startup_db_client():
    app.mongodb_client = MongoClient("mongodb://192.168.1.128:27017")
    app.database = app.mongodb_client.get_database("pymongo_tutorial")

@app.on_event("shutdown")
def shutdown_db_client():
    app.mongodb_client.close()

app.include_router(book_router, tags=["books"], prefix="/book")