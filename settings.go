package main

import (
	"net/url"
	"os"
	"strconv"
)

var DBConfig = struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}{}

// LoadConfig populates global DBConfig variable from environment
func LoadConfig() error {
	var url, err = url.Parse(os.Getenv("POSTGRES_URI"))
	if err != nil {
		return err
	}
	DBConfig.Host = url.Hostname()
	port, err := strconv.Atoi(url.Port())
	if err != nil {
		DBConfig.Port = 5432
	} else {
		DBConfig.Port = port
	}
	DBConfig.User = url.User.Username()
	password, ok := url.User.Password()
	if ok {
		DBConfig.Password = password
	}
	DBConfig.Name = url.EscapedPath()
	if len(DBConfig.Name) == 0 {
		DBConfig.Name = DBConfig.User
	} else {
		DBConfig.Name = DBConfig.Name[1:]
	}
	return nil
}
