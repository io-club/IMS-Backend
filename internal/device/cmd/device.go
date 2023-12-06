package main

import (
	"gorm.io/gorm"
	"ims-server/internal/device/dal/model"
	"ims-server/internal/user/route"
	ioconfig "ims-server/pkg/config"
	ioginx "ims-server/pkg/ginx"
	iologger "ims-server/pkg/logger"
	"ims-server/pkg/registery"
	"ims-server/pkg/util"
	"time"
)

func main() {
	config := ioconfig.GetServiceConf().Device
	addr := config.Host + ":" + config.Port
	serviceName := config.Name

	iologger.SetLogger(serviceName)

	// Get the local IP
	selfIP, err := util.GetLocalIPWithNet()
	if err != nil {
		panic(err)
	}
	// Register itself with the service registry
	heartBeat := int64(3 * time.Second)
	endpoint := selfIP + ":" + config.Port

	hub := registery.GetServiceHub(heartBeat)
	defer hub.UnRegisterService(registery.DeviceService, endpoint)

	leaseID, err := hub.RegisterService(registery.DeviceService, endpoint, 0)
	if err != nil {
		panic(err)
	}
	// Periodically send service registration
	go func() {
		for {
			hub.RegisterService(registery.DeviceService, endpoint, leaseID)
			time.Sleep(time.Duration(heartBeat) * time.Second)
		}
	}()

	// Start the service
	ioginx.NewIOServer(nil).
		InitDB(func(db *gorm.DB) error {
			return db.AutoMigrate(model.Models...)
		}).
		SetRoutes(route.Routes...).
		Run(addr, serviceName)
}
