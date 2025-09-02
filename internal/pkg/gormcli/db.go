package gormcli

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
	"wq_submitter/configs"
)

var (
	db   *gorm.DB
	once sync.Once
)

//读取配置

func openDb() {
	dbConfig := configs.GetGlobalConfig().DbConfig

	connArgs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Dbname)
	var err error
	db, err = gorm.Open(mysql.Open(connArgs), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("connect database error: %s", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("fetch database error: %s", err))
	}
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.MaxIdleTime) * time.Second)

}

func GetDb() *gorm.DB {
	once.Do(openDb)
	return db
}
