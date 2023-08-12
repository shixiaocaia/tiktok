package mysql

import (
	"fmt"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/config"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var db *gorm.DB
var dbOnce sync.Once

// OpenDB 连接mysql
func openDB() {
	// mysql连接信息
	dbConfig := config.GetGlobalConfig().MySQLConfig
	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.UserName,
		dbConfig.PassWord,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DataBase,
	)
	log.Info("mysql addr:" + connArgs)
	var err error
	db, err = gorm.Open(mysql.Open(connArgs), &gorm.Config{})
	if err != nil {
		panic("failed to connect to mysql")
	}

	// mysql连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("fetch db connection err:" + err.Error())
	}
	// 最大空闲连接数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)

	// 最大连接数
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)

	// 最大可复用时间
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.MaxIdleTime))

	log.Info("connect to mysql...")
}

// GetDB 获取db
func GetDB() *gorm.DB {
	dbOnce.Do(openDB)
	return db
}

// CloseDB 关闭db
func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			panic("fetch db connection err:" + err.Error())
		}
		sqlDB.Close()
	}
}
