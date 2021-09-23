package initialize

import (
	"mxshop-api/user-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	ApiGroup := Router.Group("/u/v1")
	router.InitUserRoutr(ApiGroup)

	return Router

}
