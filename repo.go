package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// RepoInit initializes global db variable with a DB connection and performs database migration
func RepoInit() (err error) {
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Europe/Kiev",
			DBConfig.Host,
			DBConfig.Port,
			DBConfig.User,
			DBConfig.Password,
			DBConfig.Name,
		),
	}))
	if err != nil {
		return
	}
	err = db.AutoMigrate(&Project{})
	return
}
