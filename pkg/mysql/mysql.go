package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"net/url"
	"rashomon/conf"
	"time"
)

var DB *gorm.DB // 连接主数据库的实例

func Init() {
	encoded := url.QueryEscape("Asia/Shanghai")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s", conf.AppConfig.Mysql.DbUser, conf.AppConfig.Mysql.DbPassWord, conf.AppConfig.Mysql.DbHost, conf.AppConfig.Mysql.DbPort, conf.AppConfig.Mysql.DbName, encoded)
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger:         logger,
		TranslateError: true,
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(20)  // 设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) // 打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)
}
