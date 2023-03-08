package initialize

import (
	"fish_net/cmd/api/glb"
	"fish_net/cmd/user/domain"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDB() {
	klog.Trace("connecting database...")
	db, err := gorm.Open(sqlite.Open("./cache/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.WebAuthnCredential{})
	db.AutoMigrate(&domain.Device{})
	db.AutoMigrate(&domain.Sensor{})
	db.AutoMigrate(&domain.SensorData{})

	glb.DB = db
}
