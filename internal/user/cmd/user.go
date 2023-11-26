package main

import (
	"gorm.io/gorm"
	usermodel "ims-server/internal/user/dal/model"
	"ims-server/internal/user/route"
	ioconfig "ims-server/pkg/config"
	ioginx "ims-server/pkg/ginx"
)

func main() {
	config := ioconfig.GetServiceConf().User
	addr := config.Host + ":" + config.Port
	serviceName := config.Name
	// 启动服务
	ioginx.NewIOServer(nil).
		InitDB(func(db *gorm.DB) error {
			return db.AutoMigrate(usermodel.Models...)
		}).
		SetRoutes(route.Routes...).
		Run(addr, serviceName)
}
