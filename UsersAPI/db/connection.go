package db

import (
	"UsersAPI/db/models"
	"context"
	"fmt"
	"net/url"

	_ "github.com/microsoft/go-mssqldb"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func CreateConnection(ctx context.Context) (*gorm.DB, error) {
	query := url.Values{}
	query.Add("database", "RandomUsers")

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword("sa", "Admin123"),
		Host:   fmt.Sprintf("%s:%d", "sql", 1433),
	}
	fmt.Println(u.String())
	var err error
	var database *gorm.DB
	database, err = gorm.Open(sqlserver.Open(u.String()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	_ = database.Exec(`IF NOT EXISTS(SELECT name FROM sys.databases WHERE name = 'RandomUsers')
							BEGIN
        						CREATE DATABASE RandomUsers;
							END`)

	u.RawQuery = query.Encode()
	fmt.Println(u.String())

	db, err = gorm.Open(sqlserver.Open(u.String()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Users{})

	fmt.Println("Conexion a base de datos exitosa")
	return db, nil
}

func GetDBController() *gorm.DB {
	return db
}
