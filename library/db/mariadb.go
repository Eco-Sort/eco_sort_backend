package db

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Mariadb *gorm.DB

func InitMariadb() {
	connection := os.Getenv("MARIADB_CONNECTION")
	database := os.Getenv("MARIADB_DATABASE")
	fmt.Println(connection)
	if connection == "" {
		e := errors.New("undefined MARIADB_CONNECTION")
		log.Fatal(e)
	}
	if database == "" {
		e := errors.New("undefined MARIADB_DATABASE")
		log.Fatal(e)
	}

	err := Connect(connection, database)

	if err != nil {
		log.Fatal(err)
	}
}

func Connect(connection string, database string) error {
	if Mariadb == nil {
		db, err := gorm.Open(mysql.Open(connection), &gorm.Config{})
		if err != nil {
			return err
		}
		Mariadb = db
	}
	return nil
}
