package main

import (
	"context"
	"fmt"
	"net"
	pb "taojin/grpc/proto"

	"google.golang.org/grpc"
)


type server struct{
	pb.UnimplementedUserServer
}

func (s *server) GetUser(ctx context.Context,req *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{Id: req.Id,Name: req.Name,Age: req.Age,Hobby: req.Hobby},nil
}

func main(){
	listen,_:=net.Listen("tcp","127.0.0.1:9090")

	grpcServer:=grpc.NewServer()
	pb.RegisterUserServer(grpcServer,&server{})

	err:=grpcServer.Serve(listen)
	if err!=nil{
		fmt.Printf("failed to server:%v",err)
		return
	}
}