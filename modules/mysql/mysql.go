package Mysql

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
)

type Mysql struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

func Connect() *sql.DB {
	param := &Mysql{
		Host:         "mysql-service",
		Port:         "3306",
		Username:     "user",
		Password:     "password",
		DatabaseName: "todo",
	}

	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", param.Username, param.Password, param.Host, param.Port, param.DatabaseName)

	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("Error Setup!")
		log.Fatal(err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error Ping!")
		log.Fatal(err)
	}

	return db
}
