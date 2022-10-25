from pymongo import MongoClient
from dotenv import dotenv_values
from typing import List
from pydantic import BaseModel, Field
import uuid

class Book(BaseModel):
    id: str = Field(default_factory=uuid.uuid4, alias="_id")
    title: str = Field(...)
    author: str = Field(...)
    synopsis: str = Field(...)

config = dotenv_values("../.env")

mongodb_client = MongoClient(config["ATLAS_URI"])

database = mongodb_client[config["DB_NAME"]]

books = database.get_collection("books")

resp: List[Book] = books.find(limit=100)

print(resp[0]['_id'])

mongodb_client.close()
