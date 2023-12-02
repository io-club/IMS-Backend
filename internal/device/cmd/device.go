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
	config := ioconfig.GetServiceConf().User
	addr := config.Host + ":" + config.Port
	serviceName := config.Name

	iologger.SetLogger(serviceName)

	// 获取本机 IP
	selfIP, err := util.GetLocalIPWithNet()
	if err != nil {
		panic(err)
	}
	// 向服务中心注册自己
	heartBeat := int64(3 * time.Second)
	endpoint := selfIP + ":" + config.Port

	hub := registery.GetServiceHub(heartBeat)
	defer hub.UnRegisterService(registery.DeviceService, endpoint)

	leaseID, err := hub.RegisterService(registery.DeviceService, endpoint, 0)
	if err != nil {
		panic(err)
	}
	// 周期性上报服务
	go func() {
		for {
			hub.RegisterService(registery.DeviceService, endpoint, leaseID)
			time.Sleep(time.Duration(heartBeat) * time.Second)
		}
	}()

	// 启动服务
	ioginx.NewIOServer(nil).
		InitDB(func(db *gorm.DB) error {
			return db.AutoMigrate(model.Models...)
		}).
		SetRoutes(route.Routes...).
		Run(addr, serviceName)
}
