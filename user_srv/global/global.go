package global

import (
	"log"
	"mxshop/user_srv/model"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB  *gorm.DB
	err error
)

func init() {
	dsn := "root:12345678@tcp(192.168.1.102:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // log level
			Colorful:      true,
		})
	// 全局模式

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	// 设置全局的logger,这个logger在我们执行每个sql语句的时候会打印每一行sql
	// 迁移
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
}
