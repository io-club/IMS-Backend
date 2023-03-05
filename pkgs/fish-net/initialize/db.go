package initialize

import (
	"fishnet/domain"
	"fishnet/glb"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDB() {
	glb.LOG.Debug("connecting database...")
	db, err := gorm.Open(sqlite.Open("./cache/test.db"), &gorm.Config{})
	if err != nil {
		glb.LOG.Panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.WebAuthnCredential{})
	db.AutoMigrate(&domain.Device{})
	db.AutoMigrate(&domain.Sensor{})
	db.AutoMigrate(&domain.SensorData{})

	glb.DB = db
}
