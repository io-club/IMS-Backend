package main

import (
	"ims-server/internal/device/job"
	ioconfig "ims-server/pkg/config"
	iologger "ims-server/pkg/logger"

	// "golang.org/x/vuln/client"
)

func main(){
	msg := "2:2:111"
	config := ioconfig.GetServiceConf().Device

	serviceName := config.Name
	iologger.SetLogger(serviceName)

	job.Client_test("ims" ,msg)

}