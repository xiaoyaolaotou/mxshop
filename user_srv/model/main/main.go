package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"mxshop/user_srv/model"
	"os"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	dsn := "root:12345678@tcp(192.168.1.102:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // log level
			Colorful:      true,
		})
	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}

	//options := &password.Options{16, 100, 32, sha512.New}
	//salt, encodedPwd := password.Encode("admin123", options)
	//newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	//for i := 0; i < 10; i++ {
	//	user := model.User{
	//		NickName: fmt.Sprintf("cunzhang%d", i),
	//		Mobile:   fmt.Sprintf("1803076790%d", i),
	//		Password: newPassword,
	//	}
	//	db.Save(&user)
	//}

}
