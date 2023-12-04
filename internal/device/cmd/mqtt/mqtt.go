package main

import (

	"ims-server/internal/device/job"
)

func main(){
	msg := "abc : def : hgy"
	job.ClientSub()
	job.ClientPub(msg)
}