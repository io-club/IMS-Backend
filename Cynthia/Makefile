protoc:
	protoc --go_out=./grpc/proto --go-grpc_out=./grpc/proto  ./grpc/proto/user.proto

client:
	go run ./grpc/client/client.go

server:
	go run ./grpc/server/server.go

gin:
	go run ./gin+gorm/main.go                            