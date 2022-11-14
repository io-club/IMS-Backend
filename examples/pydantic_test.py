from pydantic import BaseModel, Field
from dataclasses import dataclass
import random


# 自动生成 __init__ 方法
@dataclass
class Animal():
    name: str
    age: int


a = Animal(name="Dog", age=2)
print(a.name, a.age)


# 自动生成 __init__ 方法
def get_random_age():
    return random.randint(1, 9)


class Person(BaseModel):
    name: str
    age: int = Field(default_factory=get_random_age)

    # def __init__(self, name: str, age: int):
    #     self.name = name
    #     self.age = age


p = Person(name="aa")
print(p.name, p.age)
