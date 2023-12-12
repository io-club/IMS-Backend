package main

import (
	"gorm.io/gorm"
	"ims-server/microservices/work/internal/dal/model"
	"ims-server/microservices/work/internal/route"
	ioconfig "ims-server/pkg/config"
	ioginx "ims-server/pkg/ginx"
	iologger "ims-server/pkg/logger"
	"ims-server/pkg/registery"
	"ims-server/pkg/util"
	"time"
)

func main() {
	config := ioconfig.GetServiceConf().Work
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
	go func() {
		registery.RegisterAndReport(registery.WorkService, endpoint, heartBeat)
	}()

	// Start the service
	ioginx.NewIOServer(nil).
		InitDB(func(db *gorm.DB) error {
			return db.AutoMigrate(model.Models...)
		}).
		SetRoutes(route.Routes...).
		Run(addr, serviceName, route.Routes)
}
