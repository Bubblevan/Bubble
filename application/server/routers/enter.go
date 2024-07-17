package routers

import (
	v1 "backend/api/v1"
	"backend/global"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env)
	router := gin.Default()

	settingsApi := v1.ApiGroupApp.SettingsApi
	router.GET("", settingsApi.SettingsInfoView)

	return router
}
