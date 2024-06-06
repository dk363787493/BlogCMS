package db

import (
	"BlogCMS/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
	"os"
)

var MysqlDB *gorm.DB

func init() {
	config := config.Configuration.Mysql
	db, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 设置日志级别为 Silent，来关闭日志
	})
	if err != nil {
		slog.Error("err:", err.Error())
		os.Exit(1)
	}
	MysqlDB = db
}
