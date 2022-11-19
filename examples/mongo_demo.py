"""
mongodb 使用 demo
"""
import uuid

from pydantic import BaseModel, Field
from pymongo import MongoClient


class Book(BaseModel):
    """
    book model
    """

    id: str = Field(default_factory=uuid.uuid4, alias="_id")
    title: str = Field(...)
    author: str = Field(...)
    synopsis: str = Field(...)


mongodb_client = MongoClient("mongodb://192.168.1.128:27017")

database = mongodb_client.get_database("pymongo_tutorial")

books = database.get_collection("books")

resp = books.find(limit=100)

for book in resp:
    print(book["_id"])

mongodb_client.close()
