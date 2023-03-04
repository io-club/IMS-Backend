package glb

import (
	"fishnet/config"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	G      *gin.Engine
	Auth   *webauthn.WebAuthn
	LOG    *zap.Logger
	CONFIG *config.Config
	VP     *viper.Viper
)
