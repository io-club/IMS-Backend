from pymongo import MongoClient
from typing import List
from pydantic import BaseModel, Field
import uuid

class Book(BaseModel):
    id: str = Field(default_factory=uuid.uuid4, alias="_id")
    title: str = Field(...)
    author: str = Field(...)
    synopsis: str = Field(...)

mongodb_client = MongoClient("mongodb://192.168.1.128:27017")

database = mongodb_client.get_database("pymongo_tutorial")

books = database.get_collection("books")

resp: List[Book] = books.find(limit=100)

print(resp[0]['_id'])

mongodb_client.close()
