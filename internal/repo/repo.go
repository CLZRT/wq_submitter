package repo

import (
	"gorm.io/gorm"
	"wq_submitter/internal/pkg/gormcli"
)

var db *gorm.DB

func init() {
	db = gormcli.GetDb()
}

func GetDbCli() *gorm.DB {
	return db
}
