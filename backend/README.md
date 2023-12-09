# 实验室智能管理系统的后端仓库

本项目是 io-club 实验室智能管理系统的后端部分

## 搭建开发环境

1. 配置 git hook

   直接将相关 ./backend/.husky/hooks 下的文件复制到 ./.git/hooks/ 下即可

   `cp ./backend/.husky/hooks/* ./.git/hooks/`

2. 复制配置文件
   请按自己需要的配置修改./conf/debug.yaml，例如数据库的 ip，port 等

3. 运行

   依次运行你需要的服务，例如：`go run ./microservices/user/main.go`

   或者一键启动所有服务（需要有 make）: `make run`