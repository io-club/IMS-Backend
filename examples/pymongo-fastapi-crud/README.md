# 运行

1. 安装所需库

```sh
python -m pip install 'fastapi[all]' 'pymongo[srv]' python-dotenv
```

2. 配置 mongodb

在 _.env_ 里配置正确的信息

3. 运行

`python -m uvicorn main:app --reload`

4. 测试

`http://localhost:8000/docs`

