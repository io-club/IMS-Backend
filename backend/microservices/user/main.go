package main

import (
	"gorm.io/gorm"
	usermodel "ims-server/microservices/user/internal/dal/model"
	"ims-server/microservices/user/internal/route"
	ioconfig "ims-server/pkg/config"
	ioginx "ims-server/pkg/ginx"
	iologger "ims-server/pkg/logger"
	"ims-server/pkg/registery"
	"ims-server/pkg/util"
	"time"
)

func main() {
	config := ioconfig.GetServiceConf().User
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
		registery.RegisterAndReport(registery.UserService, endpoint, heartBeat)
	}()

	// 启动服务
	ioginx.NewIOServer(nil).
		InitDB(func(db *gorm.DB) error {
			return db.AutoMigrate(usermodel.Models...)
		}).
		SetRoutes(route.Routes...).
		Run(addr, serviceName, route.Routes)
}
