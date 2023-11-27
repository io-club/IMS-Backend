package main

import (
	"gorm.io/gorm"
	usermodel "ims-server/internal/user/dal/model"
	"ims-server/internal/user/job"
	"ims-server/internal/user/route"
	ioconfig "ims-server/pkg/config"
	ioetcd "ims-server/pkg/etcd"
	ioginx "ims-server/pkg/ginx"
	iologger "ims-server/pkg/logger"
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
	hub := job.GetServiceHub(heartBeat)
	leaseID, err := hub.RegisterService(ioetcd.UserService, selfIP+":"+config.Port, 0)
	if err != nil {
		panic(err)
	}
	// 周期性上报服务
	endpoint := selfIP + ":" + config.Port
	go func() {
		for {
			hub.RegisterService(ioetcd.UserService, endpoint, leaseID)
			time.Sleep(time.Duration(heartBeat) * time.Second)
		}
	}()

	// 启动服务
	ioginx.NewIOServer(nil).
		InitDB(func(db *gorm.DB) error {
			return db.AutoMigrate(usermodel.Models...)
		}).
		SetRoutes(route.Routes...).
		Run(addr, serviceName, hub, endpoint)
}
