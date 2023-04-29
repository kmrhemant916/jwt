package models

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
    Host     = "127.0.0.1"
    Name     = "root"
    Password = "altair"
    Database = "users"
    Port     = "3306"
)

func Connection() (*gorm.DB, error){
    connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		Name,
        Password,
        Host,
        Port,
        Database,
    )
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}