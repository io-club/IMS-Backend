package ioconfig

import (
	"github.com/spf13/viper"
	"ims-server/pkg/util"
	"log"
	"strings"
)

var (
	RootPath string // Root directory path
	V        *viper.Viper
)

func init() {
	RootPath = util.GetCurrentPath() + "/../../" // Set the root directory path

	V = NewConfigReader("debug.yaml")
}

func NewConfigReader(fileName string) *viper.Viper {
	config := viper.New()

	arr := strings.Split(fileName, ".")
	if len(arr) != 2 {
		log.Panicf("fileName must have two parts which splited by dot")
	}
	configFile, configType := arr[0], arr[1]
	config.SetConfigName(configFile)                 // Set the config file name
	config.SetConfigType(configType)                 // Set the config file type
	config.AddConfigPath(RootPath + "internal/conf") // Set the directory path where the config file is located

	// Try to find the config file
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Panicf("could not found config file %s\n", fileName)
		} else {
			log.Panicf("parse config file failed: %s\n", err.Error())
		}
	}
	return config
}
