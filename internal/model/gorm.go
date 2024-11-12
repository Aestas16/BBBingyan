package model

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "user-management-system/internal/config"
)

var db *gorm.DB

func InitDB() {
    dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", 
                       config.Config.SQL.User, config.Config.SQL.Password, config.Config.SQL.DBName, config.Config.SQL.Port)
    DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    db = DB
    if err := db.AutoMigrate(&User{}); err != nil {
        panic(err)
    }
}