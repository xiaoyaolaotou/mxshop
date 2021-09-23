package main

import (
	"mxshop-api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {
	// 1. 初始化logger
	initialize.InitLogger()
	// 2. 初始化routers
	Router := initialize.Routers()

	/*
		1. S() 可以获取一个全局的sugar, 可以让我们自己设置一个全局的logger
		2. S函数和L函数很有和， 提供了一个全局的安全访问logger途径
	*/
	zap.S().Infof("启动服务器，端口为：%d", 9999)

	Router.Run(":9999")

}
