package router

import (
	"mxshop-api/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitUserRoutr(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		UserRouter.GET("list", api.GetUserList)
	}

}
