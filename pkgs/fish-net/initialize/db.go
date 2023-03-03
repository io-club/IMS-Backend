package initialize

import (
	"demo/model"

	"demo/glb"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDB() {
	db, err := gorm.Open(sqlite.Open("./cache/test.db"), &gorm.Config{})
	if err != nil {
		glb.LOG.Panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&model.User{})

	glb.DB = db
}
