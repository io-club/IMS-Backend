package main

import (
	"context"
	"fmt"
	"log"
	pb "taojin/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func main(){
	conn,err:=grpc.Dial("127.0.0.1:9090",grpc.WithTransportCredentials(insecure.NewCredentials()))
	
	if err!=nil{
		log.Fatalf("did not connect:%v",err)
	}

	defer conn.Close()

	client:=pb.NewUserClient(conn)

	resp,err:=client.GetUser(context.Background(),&pb.UserRequest{Id: 1,Name: "Cyn",Age: 19})
	
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(resp.GetAge(),resp.GetName(),resp.GetHobby())
}