package utils

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataBaseHandler struct {
	Db  *gorm.DB
	Err error
}

type DBConfig struct {
	Name     string
	Pass     string
	Addr     string
	Database string
}

var dbHandler *DataBaseHandler

func ConDb() {
	db, err := gorm.Open(mysql.Open(os.Getenv(`DATABASE_URL`)), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting database: %v", err)
	}
	dbHandler = &DataBaseHandler{db, err}
}

func GetDb() *gorm.DB {
	return dbHandler.Db
}
